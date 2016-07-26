package db

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	// DefaultDatabaseConfigFile databases definition
	DefaultDatabaseConfigFile = "database.xml"
)

// DatabaseConfig holds general configuration of database
type DatabaseConfig struct {
	filepath  string
	xmlDbList *XMLDatabaseList
}

// SaveDatabase saves DatabaseDescriptor into the XML file
func (cfg *DatabaseConfig) SaveDatabase(database DatabaseDescriptor) {
	cfg.xmlDbList.Databases = append(cfg.xmlDbList.Databases, database)
	cfg.saveXMLFile()
}

// DropDatabase saves without database into the XML file
func (cfg *DatabaseConfig) DropDatabase(dbname string) {
	for i, v := range cfg.xmlDbList.Databases {
		if v.Name == dbname {
			cfg.xmlDbList.Databases = append(cfg.xmlDbList.Databases[:i],
				cfg.xmlDbList.Databases[i+1:]...)
			cfg.saveXMLFile()
			break
		}
	}
}

func (cfg *DatabaseConfig) saveXMLFile() {
	filePath := cfg.filepath + DefaultDatabaseConfigFile
	file, _ := os.Create(filePath)
	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(cfg.xmlDbList); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

// LoadDatabases load DatabaseConfigNode from XML file
func (cfg *DatabaseConfig) LoadDatabases() []DatabaseDescriptor {
	filePath := cfg.filepath + DefaultDatabaseConfigFile

	xmlFile, _ := os.Open(filePath)

	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)

	descriptor := XMLDatabaseList{}
	xml.Unmarshal(data, &descriptor)

	// Put the loaded database list into the sparrowdb instance list
	cfg.xmlDbList.Databases = descriptor.Databases

	v := make([]DatabaseDescriptor, 0, len(cfg.xmlDbList.Databases))

	for _, value := range cfg.xmlDbList.Databases {
		v = append(v, value)
	}

	return v
}

// NewDatabaseConfig return configuration from file
func NewDatabaseConfig(filePath string) *DatabaseConfig {
	return &DatabaseConfig{
		filepath:  filePath,
		xmlDbList: &XMLDatabaseList{},
	}
}
