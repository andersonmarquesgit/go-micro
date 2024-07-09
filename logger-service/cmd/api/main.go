package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
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
	// Conectar ao MongoDB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panicf("Can't connect to MongoDB: %v", err)
	}
	client = mongoClient

	// Criar um contexto para desconectar
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Fechar conexão
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Can't disconnect to MongoDB: %v", err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Iniciar servidor web
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
