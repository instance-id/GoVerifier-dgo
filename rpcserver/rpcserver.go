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

var (
	log *zap.SugaredLogger
	wg  sync.WaitGroup
	ad  *AccessData
)

type AccessData struct {
	key []byte
}

type ServerData struct {
	Log       *zap.SugaredLogger
	Port      *string
	Verifier  *verifier.Config
	Key       string
	Phrase    string
	encrypted []byte
}

type Server struct {
	data *ServerData
}

type Response struct {
	Message string
	Key     []byte
}

type Request struct {
	Name string
	Key  []byte
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
	if !DecryptKey(req.Key) {
		res.Message = fmt.Sprintf("Could not Authenticate")
		return nil
	}

	res.Message = fmt.Sprintf("Starting Verfier server")
	log.Infof("Starting Verifier server")
	go Run()
	return nil
}

func (s *Server) RestartServer(req Request, res *Response) error {
	if !DecryptKey(req.Key) {
		res.Message = fmt.Sprintf("Could not Authenticate")
		return nil
	}

	res.Message = fmt.Sprintf("Restarting Verfier server")
	log.Infof("Restarting Verifier server")
	runner <- true
	go Run()
	return nil
}

func (s *Server) StopServer(req Request, res *Response) error {
	if !DecryptKey(req.Key) {
		res.Message = fmt.Sprintf("Could not Authenticate")
		return nil
	}

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

func GetKey(s *ServerData) *ServerData {
	encrypted := encrypt([]byte(s.Key), s.Phrase)
	s.encrypted = encrypted
	return s
}

func DecryptKey(requestKey []byte) bool {
	unencrypted := decrypt(requestKey, serv.Phrase)
	serv.Log.Infof(fmt.Sprintf("Encrypted: %x, Decrypted: %s", requestKey, unencrypted))
	if string(unencrypted) == serv.Key {
		return true
	}
	return false
}
