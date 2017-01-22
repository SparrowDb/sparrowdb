package db

import (
	"github.com/SparrowDb/sparrowdb/db/index"
	"github.com/SparrowDb/sparrowdb/model"
	"github.com/SparrowDb/sparrowdb/util"
	"github.com/elgs/cron"
)

var (
	// keeps all active cron with dbname and id of cron
	activeCron map[string]int

	schedule *cron.Cron
)

func init() {
	activeCron = make(map[string]int)
	schedule = cron.New()
	schedule.Start()
}

func registerDbCompaction(db *Database) {
	// Cron starts compaction in another goroutine
	f, _ := schedule.AddFunc(db.Descriptor.CronExp, func() { doCompaction(db) })
	activeCron[db.Descriptor.Name] = f
}

func removeDbCompaction(dbname string) {
	if job, ok := activeCron[dbname]; ok == true {
		// removes the job from cron and delete from active cron list
		schedule.RemoveFunc(job)
		delete(activeCron, dbname)
	}
}

type tombstoneMark struct {
	path string
	index.Entry
}

func doCompaction(db *Database) {
	go db.compactionNotification()

	// get all tombstones from database
	tombstones := geTombstonesFromDb(db)

	// get all tombstones from commitlog
	tbCommitlog := getTombstonesFromCommitlog(db)
	tombstones = append(tombstones, tbCommitlog...)

	// iterate over all dataHolders
	for _, dh := range db.dhList {

		// check if DataHolder has any tombstone
		if dhContainsAnyTombstone(&dh, &tombstones) {
			// if found tombstone, get index table of DataHolder
			// and reinsert in commitlog non tombstone entry
			table := dh.summary.GetTable()

			for _, v := range table {
				if c := containsKey(v.Key, &tombstones); c == false {
					bs, _ := dh.Get(v.Offset)
					df := model.NewDataDefinitionFromByteStream(bs)
					db.commitlog.Add(df.Key, df.Status, df.Revision, bs)
				}
			}
			util.DeleteDir(dh.path)
		}
	}

	db.compFinish <- true
}

func getTombstonesFromCommitlog(db *Database) []tombstoneMark {
	var tombstones []tombstoneMark

	summary := db.commitlog.summary.GetTable()

	for _, v := range summary {
		if v.Status == model.DataDefinitionRemoved {
			tombstones = append(tombstones, tombstoneMark{db.commitlog.filepath, *v})
		}
	}
	return tombstones
}

func geTombstonesFromDb(db *Database) []tombstoneMark {
	var tombstones []tombstoneMark
	echan := make(chan []tombstoneMark)

	// iterate over all dataHolders
	for _, dh := range db.dhList {
		// search in DataHolder index for tombstone
		go func(dh *DataHolder, results chan []tombstoneMark) {
			var result []tombstoneMark

			// get index table
			idxSummary := dh.summary.GetTable()

			for _, v := range idxSummary {
				if v.Status == model.DataDefinitionRemoved {
					result = append(result, tombstoneMark{dh.path, *v})
				}
			}
			results <- result
		}(&dh, echan)
	}

	dhListLen := len(db.dhList)
	processed := 0
	for processed < dhListLen {
		select {
		case result := <-echan:
			tombstones = append(tombstones, result...)
			processed++
		}
	}

	return tombstones
}

func containsKey(key uint32, list *[]tombstoneMark) bool {
	for _, tb := range *list {
		if tb.Key == key {
			return true
		}
	}
	return false
}

func dhContainsAnyTombstone(dh *DataHolder, list *[]tombstoneMark) bool {
	var result bool
	dhs := dh.summary.GetTable()

	for _, v := range *list {
		if _, ok := dhs[v.Key]; ok == true {
			return true
		}
	}

	return result
}
