package appconfig

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

// --- Maps dbconfig.yml fields to DbSettings fields -------------------------------------------------------------------
type DbSettings struct {
	Providers []string
	Database  int `json:"database"`
	Data      struct {
		Address     string `json:"address"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		DbName      string `json:"dbname"`
		TablePrefix string `json:"tableprefix"`
	} `json:"data"`
}

// --- Gets called from Services and returns DbSettings to Dependency Injection container ------------------------------
func (d *DbSettings) GetDbConfig() *DbSettings {
	return d.loadDbConfig()
}

// --- Populates the DbSettings struct from dbconfig.yml file and returns the data for use -----------------------------
func (d *DbSettings) loadDbConfig() *DbSettings {
	config.AddDriver(yaml.Driver)
	filename := "./config/dbconfig.yml"

	err := config.LoadFiles(string(filename))
	if err != nil {
		panic(err)
	}

	dbSettings := &DbSettings{
		Providers: []string{"mysql", "postgres", "mssql", "sqlite"},
		Database:  config.Int("database"),
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
