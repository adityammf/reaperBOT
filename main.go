package main

import (
	"context"
    	"fmt"
    	"log"
    	"os"
    	"strings"
    	"time"

    	"github.com/aws/aws-lambda-go/lambda"
    	"github.com/dghubble/oauth1"
    	"github.com/reaperBOT/fetch"
    	"github.com/joho/godotenv"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context) (string, error) {
	godotenv.Load()
	
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("TOKEN_SECRET")
	prompt := os.Getenv("PROMPT")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		panic("Missing required environment variable")
	}

	fetched := fetch.GetGenerated(prompt)

	config := oath1.NewConfig(consumerKey, consumerSecret)
	token := oath1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oath1.NoContext, token)
	httpClient.Timeout = time.Second * 10

	path := "path/to/discord/server"

	body := fmt.Sprintf(`{"text": "%s"}`, fetched)

	bodyReader := strings.NewReader(body)

	response, err := 	httpClient.Post(path, "application/json", bodyReader)
	if err != nil {
		log.Fatalf("Error when posting to Discord: %v", err)
	}
	
	defer response.Body.close()
	log.Printf("Response Ok: %v", response)
	return "fin", nil
}