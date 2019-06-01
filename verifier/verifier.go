package verifier

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/instance-id/GoVerifier-dgo/appconfig"
	"github.com/instance-id/GoVerifier-dgo/verifier/cmdroutes"
)

//
//func Logger(logger *zap.Logger) {
//	discordgo.Logger = logging.DiscordgoLogger(logger.With(zap.String("feature", "discordgo")))
//}

type Verifier struct{}

type Config struct {
	Settings *appconfig.MainSettings
	di       di.Container
	Session  *discordgo.Session
	Routes   []cmdroutes.Route
	botId    string
}

func (vc *Config) registerRoutes(di di.Container, session *discordgo.Session, Routes ...cmdroutes.Route) {
	router := exrouter.New()

	cmdroutes.RegisterRoutes(
		router,
		Routes...,
	)

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		log.Printf("Message %v", m.Content)
		log.Printf("Message %v", m.Type)
		_ = router.FindAndExecute(vc.Session, vc.Settings.System.CommandPrefix, session.State.User.ID, m.Message)
	})
}

func errCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s %s\n", msg, err)
		panic(err)
	}
}

func (v *Verifier) VerifierRun(s *appconfig.MainSettings, di di.Container) (*Config, error) {
	discordgo.Logger = di.Get("logData").(func(msgL, caller int, format string, a ...interface{}))

	message := `
██╗   ██╗███████╗██████╗ ██╗███████╗██╗███████╗██████╗
██║   ██║██╔════╝██╔══██╗██║██╔════╝██║██╔════╝██╔══██╗
██║   ██║█████╗  ██████╔╝██║█████╗  ██║█████╗  ██████╔╝
╚██╗ ██╔╝██╔══╝  ██╔══██╗██║██╔══╝  ██║██╔══╝  ██╔══██╗
 ╚████╔╝ ███████╗██║  ██║██║██║     ██║███████╗██║  ██║
  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝ v0.1`
	fmt.Printf("%s\n", message)

	session, err := discordgo.New("Bot " + s.System.Token)
	if err != nil {
		fmt.Println("Error connecting to server: ", err)
		return nil, err
	}

	settings := &Config{Settings: s,
		Session: session,
		di:      di,
		Routes: []cmdroutes.Route{
			cmdroutes.NewSubRoute(di),
			cmdroutes.NewUser(di),
			cmdroutes.NewDirectMessage(),
			cmdroutes.NewAvatar(),
			cmdroutes.NewPing(),
			cmdroutes.NewListRoles(di),
			cmdroutes.NewDbSetup(di),
			cmdroutes.NewDbUpdate(di),
			cmdroutes.NewVerify(di),
			cmdroutes.NewAddUser(di)}}

	var logLevel int
	switch settings.Settings.System.ConsoleLogLevel {
	case "DEBUG":
		logLevel = discordgo.LogDebug
	case "INFO":
		logLevel = discordgo.LogInformational
	case "WARNING":
		logLevel = discordgo.LogWarning
	case "ERROR":
		logLevel = discordgo.LogError
	default:
		logLevel = discordgo.LogError
	}

	settings.Session.LogLevel = logLevel
	settings.registerRoutes(settings.di, settings.Session, settings.Routes...)

	system, err := session.User("@me")
	if err != nil {
		log.Fatalf("Unable to retrieve account information: %s ", err)
	}

	session.AddHandler(ready)

	log.Printf("Connected as: %s : %s \n", system.Username, system.Email)
	log.Printf("Command Prefix: \"%s\" registered\n", settings.Settings.System.CommandPrefix)

	return settings, nil
}

func (vc *Config) Start() error {
	err := vc.Session.Open()
	if err != nil {
		return err
	}
	return nil
}

func (vc *Config) Close() error {
	return vc.Session.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Verifier is now running. Press Ctrl + C to close session.")
}

//// Handle discord messages
//func commandHandler(config *VerifierConfig, discord *discordgo.Session, message *discordgo.MessageCreate) {
//	user := message.Author
//	if user.ID == botId || user.Bot {
//		return
//	}
//	args := strings.Split(message.Content, " ")
//	name := strings.ToLower(args[0])
//	command, found := CmdHandler.Get(name)
//	if !found {
//		return
//	}
//	channel, err := discord.State.Channel(message.ChannelID)
//	if err != nil {
//		fmt.Println("Error getting channel,", err)
//		return
//	}
//	guild, err := discord.State.Guild(channel.GuildID)
//	if err != nil {
//		fmt.Println("Error getting guild,", err)
//		return
//	}
//	ctx := CreateContext(
//		discord,
//		guild,
//		channel,
//		user,
//		message,
//		config,
//		config.CmdHandler)
//	ctx.Args = args[1:]
//	c := *command
//	c(*ctx)
//}
//
//func registerCommands(config *VerifierConfig) {
//	config.CmdHandler.Register("!v", cmd.SystemCommand)
//	config.CmdHandler.Register("!help", cmd.HelpCommand)
//}

///////////////////////////////////////////////
//router := exrouter.New()
//
//// Add some commands
//router.On("ping", func(ctx *exrouter.Context) {
//	ctx.Reply("pong")
//}).Desc("responds with pong")
//
//router.On("avatar", func(ctx *exrouter.Context) {
//	ctx.Reply(ctx.Msg.Author.AvatarURL("2048"))
//}).Desc("returns the user's avatar")
//
//// Match the regular expression user(name)?
//router.OnMatch("username", dgrouter.NewRegexMatcher("user(name)?"), func(ctx *exrouter.Context) {
//	ctx.Reply("Your username is " + ctx.Msg.Author.Username)
//})
//
//router.Default = router.On("help", func(ctx *exrouter.Context) {
//	var text = ""
//	for _, v := range router.Routes {
//		text += v.Name + " : \t" + v.Description + "\n"
//	}
//	ctx.Reply("```" + text + "```")
//}).Desc("prints this help menu")
//
//// Add message handler
//s.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
//	router.FindAndExecute(s, *fPrefix, s.State.User.ID, m.Message)
//})
//
//err = s.Open()
//if err != nil {
//log.Fatal(err)
//}
//
//log.Println("bot is running...")
//// Prevent the bot from exiting
//<-make(chan struct{})

// ---------------------------------------------------------------

//package main
//
//import (
//	"encoding/json"
//	"fmt"
//	"strings"
//
//	"github.com/andersfylling/disgord"
//	"github.com/micro/go-config"
//)
//
//var (
//	dadtype    string
//	conf_Token string
//)
//
//func init() {
//	// configuration
//	setupConf()
//}
//
//type lowconf struct {
//	Token string
//	Dad   string
//}
//
//func setupConf() {
//	_error := false
//	var conf Config
//	conf := config.NewConfig()
//	config.LoadFile("./config.json")
//	if err != nil {
//		_error = true
//		fmt.Println("Failed to read core config")
//		fmt.Println(err.Error())
//	}
//	var data lowconf
//
//	err = json.Unmarshal([]byte(file), &data)
//	if err != nil {
//		_error = true
//		fmt.Println("Failed to read core config")
//		fmt.Println(err.Error())
//	}
//
//	conf_Token = data.Token
//	dadtype = data.Dad
//
//	if !_error {
//		fmt.Println("Loaded Config")
//	} else {
//		fmt.Println("Core Config Failed to load")
//		panic(err)
//	}
//}
//
//func messageDo(session disgord.Session, data *disgord.MessageCreate) {
//	// call function to check for im responce
//	go imresponce(session, data)
//	go kysresoponce(session, data)
//	go kmsresoponce(session, data)
//
//}
//
//func kysresoponce(session disgord.Session, data *disgord.MessageCreate) {
//
//	msg := data.Message.Content
//
//	var doresponce bool
//
//	// Check if message contains kys or related
//
//	if strings.Contains(msg, " kys") {
//		doresponce = true
//	}
//	if strings.Contains(msg, " Kys") {
//		doresponce = true
//	}
//	if strings.Contains(msg, " kys ") {
//		doresponce = true
//	}
//	if strings.Contains(msg, " Kys ") {
//		doresponce = true
//	}
//	if strings.Contains(msg, "kys ") {
//		doresponce = true
//	}
//	if strings.Contains(msg, "Kys ") {
//		doresponce = true
//	}
//	if msg == "kys" {
//		doresponce = true
//	}
//	if msg == "Kys" {
//		doresponce = true
//	}
//
//	// If it does then respond
//
//	if doresponce {
//		responce := fmt.Sprint("thats not very nice, maybe you should follow your own advice before you tell someone to kys")
//		data.Message.RespondString(session, responce)
//	}
//
//}
//
//func kmsresoponce(session disgord.Session, data *disgord.MessageCreate) {
//
//	msg := data.Message.Content
//
//	var doresponce bool
//
//	// Check if message contains kms or related
//
//	if strings.Contains(msg, " kms") {
//		doresponce = true
//	}
//	if strings.Contains(msg, " Kms") {
//		doresponce = true
//	}
//	if strings.Contains(msg, " kms ") {
//		doresponce = true
//	}
//	if strings.Contains(msg, " Kms ") {
//		doresponce = true
//	}
//	if strings.Contains(msg, "kms ") {
//		doresponce = true
//	}
//	if strings.Contains(msg, "Kms ") {
//		doresponce = true
//	}
//	if msg == "kms" {
//		doresponce = true
//	}
//	if msg == "Kms" {
//		doresponce = true
//	}
//
//	// If it does then respond
//
//	if doresponce {
//		responce := fmt.Sprint("Please don't kill yourself, not under my roof, its not good for you")
//		data.Message.RespondString(session, responce)
//	}
//
//}
//
//func imresponce(session disgord.Session, data *disgord.MessageCreate) {
//
//	msg := data.Message.Content
//
//	var doresponce bool
//	var msg2 string
//
//	// Check if message starts with i'm or related
//
//	if strings.HasPrefix(msg, "i'm ") {
//		doresponce = true
//		msg2 = strings.Replace(msg, "i'm ", "", -1)
//	}
//	if strings.HasPrefix(msg, "i’m ") {
//		doresponce = true
//		msg2 = strings.Replace(msg, "i’m ", "", -1)
//	}
//	if strings.HasPrefix(msg, "im ") {
//		doresponce = true
//		msg2 = strings.Replace(msg, "im ", "", -1)
//	}
//	if strings.HasPrefix(msg, "i am ") {
//		doresponce = true
//		msg2 = strings.Replace(msg, "i am ", "", -1)
//	}
//	if strings.HasPrefix(msg, "I'm ") {
//		doresponce = true
//		msg2 = strings.Replace(msg, "I'm ", "", -1)
//	}
//	if strings.HasPrefix(msg, "I’m ") {
//		doresponce = true
//		msg2 = strings.Replace(msg, "I’m ", "", -1)
//	}
//	if strings.HasPrefix(msg, "Im ") {
//		doresponce = true
//		msg2 = strings.Replace(msg, "Im ", "", -1)
//	}
//	if strings.HasPrefix(msg, "I am ") {
//		doresponce = true
//		msg2 = strings.Replace(msg, "I am ", "", -1)
//	}
//
//	// If it does then respond with dad responce
//
//	if doresponce {
//		responce := fmt.Sprint("hi ", msg2, " i'm ", dadtype)
//		data.Message.RespondString(session, responce)
//	}
//}
//
// func main() {
// Configure disgords for it to do its shenanigans

//
//	// Message handler
//
//	session.On(disgord.EventMessageCreate, func(session disgord.Session, data *disgord.MessageCreate) {
//		fmt.Println(data.Message.Content)
//		user, err := session.GetCurrentUser()
//		if err != nil {
//			fmt.Println("Error getting current user")
//		}
//		fmt.Println(user.ID)
//		fmt.Println(data.Message.Author)
//		if data.Message.Author.ID != user.ID {
//			go messageDo(session, data)
//		}
//	})
//
//	// Create Discord Session
//
//	err = session.Connect()
//	if err != nil {
//		fmt.Println("Discord Session error")
//		fmt.Println(err.Error())
//		panic(err)
//	}
//
//	// Keep session open until ctl-c is used where it'll then be closed
//
//	session.DisconnectOnInterrupt()
//}
