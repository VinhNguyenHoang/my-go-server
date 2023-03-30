package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sendgrid/rest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	SG_API_KEY    string
	SG_PUBLIC_KEY string
	MG_DB_URI     string
	PORT          string
)

func LoadEnv() {
	env := os.Getenv("ENV")

	if env == "local" {
		fmt.Printf("LoadEnv: loading local configs...\n")
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatalf("Cannot load local env. Err: %s", err)
		}
	}

	PORT = os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("$PORT must be set")
	}

	SG_API_KEY = os.Getenv("SG_API_KEY")
	if SG_API_KEY == "" {
		log.Fatal("SendGrid API KEY not set")
	}

	SG_PUBLIC_KEY = os.Getenv("SG_PUBLIC_KEY")
	if SG_PUBLIC_KEY == "" {
		log.Fatal("SendGrid PUBLIC KEY not set")
	}

	MG_DB_URI = os.Getenv("MG_DB_URI")
	if MG_DB_URI == "" {
		log.Fatal("MongoDB URI not set")
	}
}

type EmailServer struct {
	DbConn    *mongo.Client
	GinEngine *gin.Engine
	Repo      interface {
		InsertMany(db *mongo.Client, data string)
	}

	SendGridClient interface {
		TriggerWebhookTest() (*rest.Response, error)
		SendEmail() (*rest.Response, error)
	}
}

func InitDB(ctx context.Context) *mongo.Client {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(MG_DB_URI).
		SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("error init DB: %+v", err)
	}

	return client
}

func InitEmailServer(ctx context.Context, dbConn *mongo.Client) *EmailServer {
	ginEngine := gin.Default()
	emailServer := &EmailServer{
		DbConn:         dbConn,
		GinEngine:      ginEngine,
		Repo:           &Repo{},
		SendGridClient: &SendGridClient{},
	}

	return emailServer
}

func (s *EmailServer) Run() {
	s.GinEngine.GET("/", s.Main())

	s.GinEngine.POST("/webhook", s.Webhook())

	s.GinEngine.GET("/webhooktest", s.WebhookTest())

	s.GinEngine.GET("/sendemail", s.SendEmail())

	fmt.Printf("Starting Gin server at port %s...\n", PORT)
	err := s.GinEngine.Run(":" + PORT)
	if err != nil {
		log.Fatalf("error run server: %+v", err)
	}
}

func main() {
	LoadEnv()
	ctx := context.Background()
	dbConn := InitDB(ctx)
	s := InitEmailServer(ctx, dbConn)
	s.Run()
}
