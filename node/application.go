package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/AntanasMaziliauskas/grpc/node/person"
	"github.com/phayes/freeport"
	"google.golang.org/grpc"
)

//Application struct
type Application struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	conn       *grpc.ClientConn
	Port       string
	ID         string
	ServerPort string
	lis        net.Listener
	grpcServer *grpc.Server
	Timeout    int64
	Path       string
	Person     person.PersonService
}

//Init function runs Person.Init, connects to server and sets gRPC server
func (a *Application) Init() {
	var err error

	a.wg = &sync.WaitGroup{}
	a.ctx, a.cancel = context.WithCancel(context.Background())

	if err = a.Person.Init(); err != nil {
		log.Fatalf("Init failed: %s", err)
	}

	if err = a.ConnectToServer(); err != nil {
		log.Fatalf("Did not connect to the server: %s", err)
	}

	if err = a.SetgRPCServer(); err != nil {
		log.Fatalf("Error while setting gRPC Server: %s", err)
	}

}

//Start function send a greeting to the server, launches ping go routine
//starts gRPC server
func (a *Application) Start() {

	a.ConnectWithServer()

	a.StartgRPCServer()

}

//Stop function stop gRPC server, closes connection with the server
//cancels go routines
func (a *Application) Stop() {
	a.grpcServer.Stop()
	if err := a.conn.Close(); err != nil {
		fmt.Println("Error while closing connection with the server: ", err)
	}

	a.cancel()
	a.wg.Wait()
}

//ConnectToServer function connects to server
func (a *Application) ConnectToServer() error {
	var err error

	a.conn, err = grpc.Dial(a.ServerPort, grpc.WithInsecure())

	return err
}

//ConnectWithServer function connects to a server and sends information about this Node
func (a *Application) ConnectWithServer() {
	c := api.NewNodeClient(a.conn)
	_, err := c.AddNode(context.Background(), &api.NodeInfo{Id: a.ID, Source: a.Port})
	if err != nil {
		log.Printf("Error when calling AddNode: %s", err)
	} else {
		log.Println("Connected to server.")
	}
	//	a.Timeout = response.Timeout
}

//SetgRPCServer generates random Port, sets listener, creates gRPC server object
//Attaches neccessary services to the server
func (a *Application) SetgRPCServer() error {
	var (
		err  error
		port int
	)

	if port, err = freeport.GetFreePort(); err != nil {
		log.Fatal("Failed to generate random port: ", err)
	}
	name, _ := os.Hostname()
	//a.Port = ":8080"
	//test
	a.Port = fmt.Sprintf("%s:%d", name, port)
	portas := fmt.Sprintf(":%d", port)
	if a.lis, err = net.Listen("tcp", portas); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	a.grpcServer = grpc.NewServer()
	api.RegisterServerServer(a.grpcServer, a.Person)

	return err
}

//StartgRPCServer function start gRPC server
func (a *Application) StartgRPCServer() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.grpcServer.Serve(a.lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()
}
