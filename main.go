package main

import (
	"log"
	"super-duper-fortnight/clkup"
	"super-duper-fortnight/oauth"
	"super-duper-fortnight/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	oauth.Authenticate()

	authCodeChan := make(chan string)
	go server.ServeGin(authCodeChan)
	myOAuthCode := <-authCodeChan

	token, err := clkup.GetAccessToken(myOAuthCode)
	log.Printf("OAuth token: %s", token)

}
