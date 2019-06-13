package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/instance-id/GoVerifier-dgo/appconfig"

	"github.com/sarulabs/di/v2"
	"go.uber.org/zap"

	"github.com/instance-id/GoVerifier-dgo/services"
	"github.com/instance-id/GoVerifier-dgo/verifier"
)

type appContext struct {
	Verifier *verifier.Verifier
}

var log *zap.SugaredLogger

func init() {
	message := `
██╗   ██╗███████╗██████╗ ██╗███████╗██╗███████╗██████╗
██║   ██║██╔════╝██╔══██╗██║██╔════╝██║██╔════╝██╔══██╗
██║   ██║█████╗  ██████╔╝██║█████╗  ██║█████╗  ██████╔╝
╚██╗ ██╔╝██╔══╝  ██╔══██╗██║██╔══╝  ██║██╔══╝  ██╔══██╗
 ╚████╔╝ ███████╗██║  ██║██║██║     ██║███████╗██║  ██║
  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝ v0.1`
	fmt.Printf("%s\n", message)
}

func main() {
	var appContext appContext

	log, app := DISetup()
	defer app.Delete()
	log.Infof("Initial setup complete")

	config := app.Get("configData").(*appconfig.MainSettings)
	verifierRun, err := appContext.Verifier.VerifierRun(config, app)
	ErrCheck("error creating Bot session: ", err)

	defer verifierRun.Close()
	err = verifierRun.Start()
	ErrCheck("Couldn't start verifierRun: ", err)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func DISetup() (*zap.SugaredLogger, di.Container) {
	builder, _ := di.NewBuilder()
	_ = builder.Add(services.Services...)
	app := builder.Build()
	log := app.Get("logData").(*zap.SugaredLogger)

	return log, app
}

func ErrCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s %s\n", msg, err)
	}
}
