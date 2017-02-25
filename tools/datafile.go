package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/model"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/SparrowDb/sparrowdb/util/uuid"
)

var (
	flagDataFilePath = flag.String("path", "", "Data file path (data holder or commitlog)")
)

const (
	tabSpaceSep = ' '
	version     = "1.0.0"
)

func processDataHolder(path string) {
	dataFile, err := db.OpenDataHolder(path)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	summary := dataFile.GetSummary()
	dfs := make([]*model.DataDefinition, 0)

	for _, entry := range summary.GetTable() {
		bs, err := dataFile.Get(entry.Offset)
		if err != nil {
			slog.Warnf(err.Error())
		}
		df := model.NewDataDefinitionFromByteStream(bs)
		dfs = append(dfs, df)
	}
	printTable(dfs)
}

func processCommitlog(path string) {
	path = path + ".." + string(filepath.Separator)
	cl := db.NewCommitLog(path)
	cl.LoadData()

	summary := cl.GetSummary()
	dfs := make([]*model.DataDefinition, 0)

	for _, entry := range summary.GetTable() {
		bs := cl.GetByHash(entry.Key)
		df := model.NewDataDefinitionFromByteStream(bs)
		dfs = append(dfs, df)
	}
	printTable(dfs)
}

func printTable(dfs []*model.DataDefinition) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, tabSpaceSep, tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%v\t%s\t%s\t%s", "Key", "Ext", "Size", "Status", "Revision", "Timestamp"))
	for _, df := range dfs {
		uuid, err := uuid.ParseUUID(df.Token)
		if err != nil {
			slog.Warnf(err.Error())
			continue
		}
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%v\t%v\t%v\t%s", df.Key, df.Ext, df.Size, df.Status, df.Revision, uuid.Time().String()))
	}
	w.Flush()

}

func main() {
	flag.Parse()
	slog.Infof("SparrowDb Tool %s - data visualizer", version)

	if *flagDataFilePath == "" {
		slog.Fatalf("Invalid data path")
	}

	dirInfo, err := os.Stat(*flagDataFilePath)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	abspath, err := filepath.Abs(*flagDataFilePath)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	slog.Infof("Data file: %s", abspath)

	if dirInfo.Name() == "commitlog" {
		processCommitlog(*flagDataFilePath)
	} else {
		processDataHolder(*flagDataFilePath)
	}
}
