package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/instance-id/GoVerifier-dgo/appconfig"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"go.uber.org/zap"

	"github.com/instance-id/GoVerifier-dgo/logging"
	"github.com/instance-id/GoVerifier-dgo/services"
	"github.com/instance-id/GoVerifier-dgo/verifier"
)

type appContext struct {
	Verifier *verifier.Verifier
}

var (
	service = "Verifier"
)

func main() {
	var appContext appContext

	log, err := logging.NewLogger(
		logging.DevelopmentEnvironment,
		service,
		"",
		&http.Client{
			Timeout: 10 * time.Second,
		})

	func(log *zap.Logger) {
		discordgo.Logger = logging.DiscordgoLogger(log.With(zap.String("feature", "discordgo")))
	}(log)

	builder, err := di.NewBuilder()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = builder.Add(services.Services...)
	if err != nil {
		log.Fatal("Error", zap.Error(err))
	}
	app := builder.Build()
	defer app.Delete()

	guildObject, err := app.SafeGet("configData")
	if guild, ok := guildObject.(*appconfig.MainSettings); ok {
		log.Sugar().Infof("GuildID: %s", guild.Discord.Guild)
	} else {
		log.Sugar().Infof("Token: %s", guild.Discord.Guild)
	}

	verifierRun, err := appContext.Verifier.VerifierRun(app.Get("configData").(*appconfig.MainSettings), app)
	if err != nil {
		log.Sugar().Fatalf("error creating Bot session,", err)
	}

	defer verifierRun.Close()
	err = verifierRun.Start()
	if err != nil {
		log.Sugar().Fatalf("Couldn't start verifierRun: %v", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
