package verifier

import (
	. "github.com/instance-id/GoVerifier-dgo/utils"
	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/instance-id/GoVerifier-dgo/appconfig"
	"github.com/instance-id/GoVerifier-dgo/verifier/cmdroutes"
)

type Verifier struct{}

var log *zap.SugaredLogger

type Config struct {
	Settings *appconfig.MainSettings
	di       di.Container
	Session  *discordgo.Session
	Routes   []cmdroutes.Route
	BotId    string
}

func (v *Verifier) VerifierRun(s *appconfig.MainSettings, di di.Container) (*Config, error) {
	log = di.Get("logData").(*zap.SugaredLogger)
	CmdInitialize(di)

	session, err := discordgo.New("Bot " + s.System.Token)
	ErrCheck("Error connecting to server: ", err)

	settings := &Config{Settings: s,
		Session: session,
		di:      di,
		Routes: []cmdroutes.Route{
			cmdroutes.NewSubRoute(di),
			cmdroutes.NewUser(di),
			cmdroutes.NewDirectMessage(),
			cmdroutes.NewListRoles(di),
			cmdroutes.NewDbSetup(di),
			cmdroutes.NewDbUpdate(di),
			cmdroutes.NewVerify(di),
			cmdroutes.NewAddUser(di),
			cmdroutes.NewRequest(di),
			cmdroutes.NewTest(di)}}

	var logLevel int
	switch settings.Settings.System.FileLogLevel {
	case 0:
		logLevel = discordgo.LogInformational
	case 1:
		logLevel = discordgo.LogDebug
	case 2:
		logLevel = discordgo.LogWarning
	case 3:
		logLevel = discordgo.LogError
	default:
		logLevel = discordgo.LogInformational
	}

	settings.Session.LogLevel = logLevel
	settings.registerRoutes(settings.di, settings.Session, settings.Routes...)

	system, err := session.User("@me")
	ErrCheck("Unable to retrieve account information: ", err)

	session.AddHandler(ready)

	log.Infof("Connected as: %s : %s ", system.Username, system.Email)
	log.Infof("Command Prefix: \"%s\" registered", settings.Settings.System.CommandPrefix)

	return settings, nil
}

func (vc *Config) Start() error {
	err := vc.Session.Open()
	ErrCheck("Could not start server: ", err)
	return nil
}

func (vc *Config) Close() error {
	return vc.Session.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Infof("Verifier is now running. Press Ctrl + C to close session.")
}

func (vc *Config) registerRoutes(di di.Container, session *discordgo.Session, Routes ...cmdroutes.Route) {
	router := exrouter.New()

	cmdroutes.RegisterRoutes(
		router,
		Routes...,
	)

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		log.Infof("Message %v", m.Content)
		log.Infof("Message %v", m.Type)
		_ = router.FindAndExecute(vc.Session, vc.Settings.System.CommandPrefix, session.State.User.ID, m.Message)
	})
}

func ErrCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s %s\n", msg, err)
	}
}
