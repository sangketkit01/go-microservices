package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/sangketkit01/logger-service/data"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx) ; err != nil{
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(mongoClient),
	}

	// register the RPC server
	err = rpc.Register(new(RPCServer))

	// start rpc server
	go app.rpcListen()

	// start grpc server
	go app.gRPCListen()
	
	// start web server 
	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
	}
}

func(app *Config) rpcListen()  error {
	log.Println("Starting RPC server on port:", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil{
		return err
	}

	defer listen.Close()

	for{
		rpcConn, err := listen.Accept()
		if err != nil{
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	conn, err := mongo.Connect(clientOptions)
	if err != nil{
		log.Println("Error connecting:", err)
		return nil, err
	}

	return conn, nil
}
