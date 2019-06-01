package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/instance-id/GoVerifier-dgo/components"

	"github.com/instance-id/GoVerifier-dgo/logging"

	"github.com/bwmarrin/discordgo"
	"github.com/instance-id/GoVerifier-dgo/appconfig"
	"github.com/sarulabs/di/v2"
	"go.uber.org/zap"
)

var Services = []di.Def{
	{
		Name:  "configData",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			var cfg appconfig.MainSettings
			config := cfg.GetConfig()
			return config, nil
		}},
	{
		// --- Creates database connection object ----------------------------------------------------------------------
		Name:  "dbConn",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			var db appconfig.DbSettings
			var conn components.DbConfig
			dbconfig := db.GetDbConfig()
			dbConn := conn.ConnectDB(dbconfig)
			return dbConn, nil
		},
		Close: func(obj interface{}) error {
			return obj.(*components.DbConfig).Xorm.Close()
		},
	},
	{
		// --- Uses database connection object and returns a connection session ----------------------------------------
		Name:  "db",
		Scope: di.Request,
		Build: func(ctn di.Container) (interface{}, error) {
			conn := ctn.Get("dbConn").(*components.DbConfig).Xorm
			return conn, nil
		},
		Close: func(obj interface{}) error {
			return obj.(*components.DbConfig).Xorm.Close()
		},
	},
	{
		// ---Creates Zap to default DiscordGo logger ------------------------------------------------------------------
		Name: "logData",
		Build: func(ctn di.Container) (interface{}, error) {
			var service = "Verifier"

			logger, err := logging.NewLogger(
				logging.DevelopmentEnvironment,
				service,
				"",
				&http.Client{
					Timeout: 10 * time.Second,
				})
			if err != nil {
				fmt.Printf("Could not get new logger! %s", err)
			}

			f := func(logger *zap.Logger) func(msgL, caller int, format string, a ...interface{}) {
				discordgo.Logger = logging.DiscordgoLogger(logger.With(zap.String("feature", "discordgo")))
				return discordgo.Logger
			}
			return f(logger), nil
		}},
}
