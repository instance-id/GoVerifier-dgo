package rpcserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/instance-id/GoVerifier-dgo/verifier"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap"
)

var (
	serv       *RpcServer
	log        *zap.SugaredLogger
	wg         sync.WaitGroup
	ad         *AccessData
	runner     chan bool
	timeOut    chan bool
	ip         = "localhost"
	port       = "14555"
	address    = fmt.Sprintf("%s:%s", ip, port)
	clientConn net.Conn
)

type Server struct{}

type AccessData struct {
	key []byte
}

type ServerData struct {
	ProcessName     string
	ProcessID       int
	Log             *zap.SugaredLogger
	Address         string
	Port            *string
	Verifier        *verifier.Config
	Key             string
	Phrase          string
	encrypted       []byte
	RpcRunning      bool
	VerifierRunning bool
	svr             *server.Server
}

type RpcServer struct {
	Data *ServerData
}

type Reply struct {
	Message   string
	RunCheck  bool
	RPCCheck  bool
	ProcessID int
	Key       []byte
}

type Args struct {
	Key     []byte
	Name    string
	Message string
	Ctx     map[string]interface{}
}

func Status(reply *Reply) *Reply {
	log.Infof("Checking RPC_STATUS")

	if serv.Data.VerifierRunning {
		log.Infof("Verifier UI checking status")
		reply.Message = fmt.Sprintf("Verifier running")
		reply.RunCheck = serv.Data.VerifierRunning
		reply.RPCCheck = serv.Data.RpcRunning
		reply.ProcessID = os.Getpid()

		return reply

	} else {
		log.Infof("Verifier not running")
		reply.Message = "Verifier not running"
		reply.RunCheck = false
		reply.RPCCheck = true

		return reply
	}
}

func StartServer(args Args, res *Reply) *Reply {
	res.Message = fmt.Sprintf("Starting Verfier server")
	go Run()

	log.Infof("Starting Verifier server")
	return res
}

func RestartServer(args Args, res *Reply) *Reply {

	log.Infof("Stopping Verifier server")
	serv.Data.Verifier.Close()
	serv.Data.VerifierRunning = false
	log.Infof("Restarting Verifier server")

	defer serv.Data.Verifier.Close()
	err := serv.Data.Verifier.Start()
	ErrCheck("Couldn't start verifierRun: ", err)
	serv.Data.VerifierRunning = true

	log.Infof("Restart complete")
	res.Message = fmt.Sprintf("Restart complete")

	return res
}

func StopServer(args Args, res *Reply) *Reply {
	res.Message = fmt.Sprintf("Shut down of Verifier initiated")
	res.RunCheck = false
	res.RPCCheck = false
	go delayedOsExit()
	return res
}

func RunServer(r *RpcServer, Logs *zap.SugaredLogger) {
	serv = r
	log = Logs
	serv.Data.VerifierRunning = false
	wg.Add(1)

	go rpcServer()
	runtime.Gosched()
	log.Infof("RPC Server started. Waiting for Verifier start signal.")
	log.Infof("Initiated 30 second application timeout if not received.")

	serv.Data.RpcRunning = true
	receiveOrTimeout()
	wg.Wait()
	log.Infof("Shutdown Initiated")
	serv.Data.VerifierRunning = false
	serv.Data.RpcRunning = false

}

func rpcServer() {
	log.Infof("Starting Server on port: %s", RPC_PORT)

	s := server.NewServer()
	//s.Register(new(Arith), "")
	s.RegisterName("Server", new(Server), "")
	err := s.Serve("tcp", address)
	if err != nil {
		panic(err)
	}
}

func Run() {

	defer serv.Data.Verifier.Close()
	err := serv.Data.Verifier.Start()
	ErrCheck("Couldn't start verifierRun: ", err)
	serv.Data.VerifierRunning = true
	runner = make(chan bool)
	timeOut <- true
	<-runner
	wg.Done()
	log.Infof("Issued close")
}

func ErrCheck(msg string, err error) {
	if err != nil {
		log.Fatalf("%s %s\n", msg, err)
	}
}

func GetKey(s *ServerData) *ServerData {
	encrypted := encrypt([]byte(s.Key), s.Phrase)
	s.encrypted = encrypted
	return s
}

func DecryptKey(requestKey []byte) bool {
	unencrypted := decrypt(requestKey, serv.Data.Phrase)
	if string(unencrypted) == serv.Data.Key {
		return true
	}
	return false
}

func delayedOsExit() {
	time.Sleep(5000 * time.Millisecond)
	log.Infof("Stopping Verifier server")
	runner <- true
}

func (s *Server) RpcRequestHandler(ctx context.Context, args Args, reply *Reply) error {
	fmt.Printf("%x : %s : %s: %s \n", args.Key, args.Name, args.Message, reply.Message)

	clientConn = ctx.Value(server.RemoteConnContextKey).(net.Conn)
	log.Infof(" Client IP: %s ", clientConn.RemoteAddr().String())

	if !DecryptKey(args.Key) {
		log.Infof("Failed Authentication attempt from: %s ", clientConn.RemoteAddr().String())
		reply.Message = fmt.Sprintf("Could not Authenticate from: %s", clientConn.RemoteAddr().String())
		return nil
	}
	serv.Data.Log.Infof("Authentication successful")

	switch args.Name {
	case RPC_START:
		{
			reply = StartServer(args, reply)
			reply.ProcessID = os.Getpid()
			return nil
		}
	case RPC_RESTART:
		{
			reply = RestartServer(args, reply)
			return nil
		}
	case RPC_STOP:
		{
			reply = StopServer(args, reply)
			return nil
		}
	case RPC_STATUS:
		{
			log.Infof("Received RPC_STATUS cmd")
			reply = Status(reply)
			log.Infof("Finished Status. \n")
			return nil
		}
	default:
		{
			log.Infof("None matched up. : / \n")
			return nil
		}
	}
}

func receiveOrTimeout() {

	timeOut = make(chan bool, 1)
	defer close(timeOut)

	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()

	select {
	case <-timeOut:
		log.Infof("Start signal received. Timeout canceled.")
	case <-timer.C:
		log.Warnf("Start signal not received before timeout. Closing.")
		wg.Done()
	}
}
