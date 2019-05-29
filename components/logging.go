package components

import (
	"os"

	"github.com/kz/discordrus"

	"fmt"
	"runtime"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/instance-id/GoVerifier/cache"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// InitLogger initialises and caches the logging server
// will send warning and above messages to Discord if the DISCORD_LOGGING_WEBHOOK_URL environment variable is set
func InitLogger(service string) {
	format := new(prefixed.TextFormatter)
	format.TimestampFormat = "02-01-06 15:04:05.000"
	format.FullTimestamp = true
	format.ForceColors = true
	format.SpacePadding = 2

	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	log.Formatter = format

	// log.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "02-01-06 15:04:05.000"}
	log.Hooks = make(logrus.LevelHooks)

	// send warnings and above to discord if DISCORD_LOGGING_WEBHOOK_URL is set
	if os.Getenv("DISCORD_LOGGING_WEBHOOK_URL") != "" {
		log.AddHook(discordrus.NewHook(
			os.Getenv("DISCORD_LOGGING_WEBHOOK_URL"),
			logrus.WarnLevel,
			&discordrus.Opts{
				Username:         "Logging",
				Author:           "at " + service,
				DisableTimestamp: true,
			},
		))
	}

	// setup discord logging
	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		pc, file, line, _ := runtime.Caller(caller)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := format
		if strings.Contains(msg, "%") {
			msg = fmt.Sprintf(format, a...)
		}

		switch msgL {
		case discordgo.LogError:
			log.WithFields(logrus.Fields{"service": service, "module": "discordgo"}).
				Errorf("%s:%d:%s() %s", file, line, name, msg)
		case discordgo.LogWarning:
			log.WithFields(logrus.Fields{"service": service, "module": "discordgo"}).
				Warnf("%s:%d:%s() %s", file, line, name, msg)
		case discordgo.LogInformational:
			log.WithFields(logrus.Fields{"service": service, "module": "discordgo"}).
				Infof("%s:%d:%s() %s", file, line, name, msg)
		case discordgo.LogDebug:
			log.WithFields(logrus.Fields{"service": service, "module": "discordgo"}).
				Debugf("%s:%d:%s() %s", file, line, name, msg)
		}
	}

	cache.SetLogger(log.WithFields(logrus.Fields{"service": service}))
}
