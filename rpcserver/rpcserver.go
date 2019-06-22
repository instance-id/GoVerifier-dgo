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

func Status(args Args, res *Reply) *Reply {
	var tmpReply = new(Reply)
	log.Infof("Checking RPC_STATUS")

	if serv.Data.VerifierRunning {
		log.Infof("Verifier UI checking status")
		tmpReply.Message = fmt.Sprintf("Starting Verfier server")
		tmpReply.RunCheck = serv.Data.VerifierRunning
		tmpReply.RPCCheck = serv.Data.RpcRunning
		return tmpReply

		//res.ProcessID = serv.Data.svr.
	} else {
		log.Infof("Verifier not running")
		tmpReply.Message = "Verifier not running"
		tmpReply.RunCheck = false
		tmpReply.RPCCheck = true
		return tmpReply
	}
}

func StartServer(args Args, res *Reply) *Reply {
	res.Message = fmt.Sprintf("Starting Verfier server")
	go Run()

	log.Infof("Starting Verifier server: After Run()")
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
	serv.Data.RpcRunning = true
	receiveOrTimeout()
	wg.Wait()
	log.Infof("Canceled via RPC")
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
	var tmpReply *Reply

	clientConn = ctx.Value(server.RemoteConnContextKey).(net.Conn)
	log.Infof(" Client IP: %s \n ", clientConn.RemoteAddr().String())

	if !DecryptKey(args.Key) {
		log.Warnf("Failed Authentication attempt from: %s ", clientConn.RemoteAddr().String())
		args.Message = fmt.Sprintf("Could not Authenticate from: %s", clientConn.RemoteAddr().String())
		return nil
	}
	serv.Data.Log.Infof("Authentication successful")

	switch args.Name {
	case RPC_START:
		{
			//func() *Reply { tmpReply = StartServer(args, reply); return tmpReply }()
			tmpReply = StartServer(args, reply)
			tmpReply.ProcessID = os.Getpid()
		}
	case RPC_RESTART:
		{
			func() *Reply { tmpReply = RestartServer(args, reply); return tmpReply }()
		}
	case RPC_STOP:
		{
			func() *Reply { tmpReply = StopServer(args, reply); return tmpReply }()
		}
	case RPC_STATUS:
		{
			log.Infof("Received RPC_STATUS cmd")
			//var answer string
			//answer += "Process name: " + serv.Data.ProcessName + "\n"

			//func() *Reply { tmpReply = Status(args, reply); return tmpReply }()
			tmpReply = Status(args, reply)
			log.Infof("Finished Status.\n")

		}
	default:
		{
			log.Infof("None matched up. : / \n")
		}
	}

	log.Infof("%x : %s : %s : %s", args.Key, args.Name, args.Message, reply.Message)
	reply = tmpReply
	return nil
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
