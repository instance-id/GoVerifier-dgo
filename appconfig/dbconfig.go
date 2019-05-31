package appconfig

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type DbData struct {
	DbSettings DbSettings
}

// --- Maps dbconfig.yml fields to DbSettings fields -------------------------------------------------------------------
type DbSettings struct {
	Database string `json:"database"`
	Data     struct {
		Address     string `json:"address"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		DbName      string `json:"dbname"`
		TablePrefix string `json:"tableprefix"`
	} `json:"data"`
}

// --- Gets called from Services and returns DbSettings to Dependency Injection container ------------------------------
func (d *DbSettings) GetDbConfig() *DbSettings {
	return d.LoadDbConfig()
}

// --- Populates the DbSettings struct from dbconfig.yml file and returns the data for use -----------------------------
func (d *DbSettings) LoadDbConfig() *DbSettings {
	config.AddDriver(yaml.Driver)
	filename := "./appconfig/dbconfig.yml"

	err := config.LoadFiles(string(filename))
	if err != nil {
		panic(err)
	}

	dbSettings := &DbSettings{
		Database: config.String("database"),
		Data: struct {
			Address     string `json:"address"`
			Username    string `json:"username"`
			Password    string `json:"password"`
			DbName      string `json:"dbname"`
			TablePrefix string `json:"tableprefix"`
		}{
			Address:     config.String("data.address"),
			Username:    config.String("data.username"),
			Password:    config.String("data.password"),
			DbName:      config.String("data.dbname"),
			TablePrefix: config.String("data.tableprefix"),
		},
	}
	return dbSettings
}
