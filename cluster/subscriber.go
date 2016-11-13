package cluster

import (
	"encoding/json"

	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/errors"
	"github.com/SparrowDb/sparrowdb/model"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/SparrowDb/sparrowdb/spql"
	"github.com/SparrowDb/sparrowdb/util"
	"github.com/SparrowDb/sparrowdb/util/uuid"
	"github.com/nats-io/nats"
)

var (
	_config        *db.SparrowConfig
	_dbm           *db.DBManager
	_queryExecutor *spql.QueryExecutor
	_connection    *nats.Conn
	_encon         *nats.EncodedConn
	_enconData     *nats.EncodedConn

	_sparrowQuerySub = "sparrow.query"
	_sparrowDataSub  = "sparrow.data"

	chRecvData chan string
	chSendData chan string
)

type message struct {
	Name    string
	Vars    map[string]string
	Content interface{}
}

func connect() {
	var err error
	_connection, err = nats.Connect(_config.PublisherServers)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	_encon, err = nats.NewEncodedConn(_connection, nats.JSON_ENCODER)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	_enconData, err = nats.NewEncodedConn(_connection, nats.GOB_ENCODER)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	registerReceiverBinder()
}

func registerReceiverBinder() {
	_encon.Subscribe(_sparrowQuerySub, func(m *message) {
		if m.Name != _config.NodeName {
			var qr spql.QueryRequest
			str := m.Content.(string)
			json.Unmarshal([]byte(str), &qr)

			q, _ := qr.ParseQuery()
			results := <-_queryExecutor.ExecuteQuery(&q)
			if results == nil {
				slog.Errorf(errors.ErrEmptyQueryResult.Error())
			}
		}
	})

	_enconData.Subscribe(_sparrowDataSub, func(m *message) {
		if m.Name != _config.NodeName {
			dbname := m.Vars["database"]

			if db, ok := _dbm.GetDatabase(dbname); ok == true {
				bs := util.NewByteStreamFromBytes(m.Content.([]byte))
				df := model.NewDataDefinitionFromByteStream(bs)

				storedDf, found := db.GetDataByKey(df.Key)

				if found {
					tm, _ := uuid.ParseUUID(df.Token)
					stm, _ := uuid.ParseUUID(storedDf.Token)

					if tm.Time().Before(stm.Time()) || tm.Time().Equal(stm.Time()) {
						return
					}

					df.Revision = storedDf.Revision
					df.Revision++
				}

				db.InsertCheckUpsert(df, true)
			}
		}
	})
}

// PublishQuery publishes query
func PublishQuery(query spql.QueryRequest) {
	if query.Action == "select" {
		return
	}

	b, _ := json.Marshal(query)
	m := message{
		_config.NodeName,
		map[string]string{},
		string(b),
	}
	_encon.Publish(_sparrowQuerySub, m)
}

// PublishData plushes data
func PublishData(df model.DataDefinition, dbname string) {
	m := message{
		_config.NodeName,
		map[string]string{
			"database": dbname,
		},
		df.ToByteStream().Bytes(),
	}
	_enconData.Publish(_sparrowDataSub, m)
}

// Close finishes cluster service
func Close() {
	slog.Infof("Stopping Cluster service")
	_encon.Close()
	_enconData.Close()
	_connection.Close()
}

// Start Starts cluster service
func Start(config *db.SparrowConfig, dbm *db.DBManager) {
	slog.Infof("Starting Cluster service")
	_config = config
	_dbm = dbm
	_queryExecutor = spql.NewQueryExecutor(dbm)
	connect()
}
