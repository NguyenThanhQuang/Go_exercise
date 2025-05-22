package main

import (
	"os"
	"social_media_server/config"
	"social_media_server/routes"
)

func main() {
	config.ConnectDatabase()
	r := routes.SetupRouter()
	r.Run(os.Getenv("PORT"))
}