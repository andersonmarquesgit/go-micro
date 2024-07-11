package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// Connect MongoDB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panicf("Can't connect to MongoDB: %v", err)
	}
	client = mongoClient

	// Create context for disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Can't disconnect to MongoDB: %v", err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Register the RPC Server
	err = rpc.Register(new(RPCServer))

	go app.rpcListen()

	// start wev server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting logger service on port %s\n", webPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Can't start logger service %v", err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}

}

func connectToMongo() (*mongo.Client, error) {
	// Obter URL do MongoDB das variáveis de ambiente
	mongoURL := os.Getenv("MONGO_URI")
	if mongoURL == "" {
		mongoURL = "mongodb://admin:password@mongo:27017"
	}

	// Criar opções de conexão
	clientOptions := options.Client().ApplyURI(mongoURL)

	// Conectar ao MongoDB
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Can't connect to MongoDB: ", err)
		return nil, err
	}

	// Verificar conexão
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Can't ping to MongoDB: ", err)
		return nil, err
	}

	log.Println("Success MongoDB connect")
	return c, nil
}
