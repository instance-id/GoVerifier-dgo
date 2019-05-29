package services

import (
	"fmt"
	"net/http"
	"time"

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
		}}, {
		Name: "logData",
		Build: func(ctn di.Container) (interface{}, error) {
			var (
				service = "Verifier"
			)

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
