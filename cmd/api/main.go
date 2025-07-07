package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"pawrest/internal/api/middleware"
	"pawrest/internal/api/routes"
	"pawrest/internal/cfgyaml"
	"pawrest/internal/db"
)

func main() {
	httpsFlag := flag.Bool("https", false, "Start the server with HTTPS")
	portFlag := flag.String("port", "", "Server port")
	certFlag := flag.String("cert", "keys/server.pem", "TLS certificate file location")
	keyFlag := flag.String("key", "keys/server.key", "TLS private key file location")
	flag.Parse()

	if err := cfgyaml.Load("env.yaml"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to the database...")
	database, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	if err := middleware.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logging middleware: %v\n", err)
	}
	defer middleware.CloseLogger()

	fmt.Println("Starting up the server...")
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(middleware.FileLogger())
	routes.Router(router, database)

	useHTTPS := *httpsFlag
	port := *portFlag

	if port == "" {
		if useHTTPS {
			port = "8443"
		} else {
			port = "8080"
		}
	}

	if useHTTPS {
		if err := router.RunTLS(":"+port, *certFlag, *keyFlag); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := router.Run(":" + port); err != nil {
			log.Fatal(err)
		}
	}
}
