package main

import (
	"os"

	"example-loan.com/loan-app/app/config"
	"example-loan.com/loan-app/app/helpers"
	"example-loan.com/loan-app/app/router"
	"github.com/joho/godotenv"
)

func main() {
	branch := helpers.GetGitBranch()
	envFile := helpers.GetEnvFileFromBranch(branch)
	err := godotenv.Load(envFile)
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	initApp, err := config.InitializeApp()
	if err != nil {
		panic(err)
	}

	router := router.InitRouter(initApp)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
