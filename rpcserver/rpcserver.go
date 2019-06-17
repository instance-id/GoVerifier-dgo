package rpcserver

import (
	"fmt"
	"net"
	"net/rpc"
	"runtime"
	"sync"

	"github.com/instance-id/GoVerifier-dgo/verifier"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger

var wg sync.WaitGroup

type ServerData struct {
	Log      *zap.SugaredLogger
	Port     *string
	Verifier *verifier.Config
}

type Server struct{}

type Response struct {
	Message string
}

type Request struct {
	Name string
}

var (
	runner chan bool
	serv   *ServerData
)

func (s *Server) Add(u [2]int64, reply *int64) error {
	*reply = u[0] + u[1]
	fmt.Println("Received connection. Executing command!")
	return nil
}

func (s *Server) StartServer(req Request, res *Response) error {
	res.Message = fmt.Sprintf("Starting Verfier server")
	log.Infof("Starting Verifier server")
	go Run()
	return nil
}

func (s *Server) RestartServer(req Request, res *Response) error {
	res.Message = fmt.Sprintf("Restarting Verfier server")
	log.Infof("Restarting Verifier server")
	runner <- true
	go Run()
	return nil
}

func (s *Server) StopServer(req Request, res *Response) error {
	res.Message = fmt.Sprintf("Stopping Verfier server")
	log.Infof("Stopping Verifier server")
	runner <- true
	return nil
}

func server() {
	log.Infof("Starting Server!")
	log.Infof(fmt.Sprintf("127.0.0.1:%s", *serv.Port))

	err := rpc.Register(new(Server))
	if err != nil {
		fmt.Printf("Error registering RPC: %s", err)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", *serv.Port))
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Infof("RPC server running at: %s", ln.Addr().String())

	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(c)
	}
}

func RunServer(s *ServerData, logs *zap.SugaredLogger) {
	wg.Add(1)
	log = logs
	serv = s
	go server()
	runtime.Gosched()
	log.Infof("RPC Server started. Waiting for Verifier start signal.")
	wg.Wait()
	log.Infof("Canceled via RPC")

}

func Run() {
	defer serv.Verifier.Close()
	err := serv.Verifier.Start()
	ErrCheck("Couldn't start verifierRun: ", err)

	runner = make(chan bool)
	<-runner
	wg.Done()
	log.Infof("Issued close")

}

func ErrCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s %s\n", msg, err)
	}
}

//func Timeout() *time.Timer {
//	timer1 := time.NewTimer(10 * time.Second)
//	<-timer1.C
//	<-runner
//	fmt.Println("Timer 1 expired")
//	return timer1
//}
