package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/instance-id/GoVerifier-dgo/appconfig"
	. "github.com/instance-id/GoVerifier-dgo/rpcserver"

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

}

func main() {
	var appContext appContext
	var port *string

	useRPC := flag.Bool("rpc", false, "Run RPC server for UI communication?")
	port = flag.String("port", "14555", "Port for RPC commands")
	flag.Parse()

	if !*useRPC {
		outPath := "./logs/verifier.log"
		outPathLast := "./logs/verifierlastrun.log"
		if err := os.Rename(outPath, outPathLast); err != nil {
			log.Infof("Could not rotate log file")
		}
	}

	log, app := DISetup()
	defer app.Delete()

	message := []string{
		"██╗   ██╗███████╗██████╗ ██╗███████╗██╗███████╗██████╗",
		"██║   ██║██╔════╝██╔══██╗██║██╔════╝██║██╔════╝██╔══██╗",
		"██║   ██║█████╗  ██████╔╝██║█████╗  ██║█████╗  ██████╔╝",
		"╚██╗ ██╔╝██╔══╝  ██╔══██╗██║██╔══╝  ██║██╔══╝  ██╔══██╗",
		" ╚████╔╝ ███████╗██║  ██║██║██║     ██║███████╗██║  ██║",
		"  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝ v0.1"}
	for s := range message {
		msg := fmt.Sprintf("%s", message[s])
		log.Infof("%s", msg)
	}

	config := app.Get("configData").(*appconfig.MainSettings)
	verifierRun, err := appContext.Verifier.VerifierRun(config, app)
	ErrCheck("error creating Bot session: ", err)

	if *useRPC {

		log.Infof("Configuration import complete. Starting Verifier RPC server.")
		var rpcServer = RpcServer{
			Data: &ServerData{
				Port:     port,
				Log:      log,
				Verifier: verifierRun,
				Phrase:   config.Discord.Guild,
				Key:      config.System.Token[len(config.System.Token)-13:],
			},
		}

		RunServer(&rpcServer, log)

	} else {

		log.Infof("Configuration import complete. Starting Verifier in standalone mode.")

		defer verifierRun.Close()
		err = verifierRun.Start()
		ErrCheck("Couldn't start verifierRun: ", err)

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}
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
