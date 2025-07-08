package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"pawrest/internal/api/middleware"
	"pawrest/internal/api/routes"
	"pawrest/internal/cfgyaml"
	"pawrest/internal/db"
)

func isFlagPassed(flagName string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			found = true
		}
	})

	return found
}

func resolveStrFlag(flagPtr *string, flagName, envVar string) string {
	if isFlagPassed(flagName) {
		return *flagPtr
	}

	if val := os.Getenv(envVar); val != "" {
		return val
	}

	return *flagPtr
}

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

	useHTTPS := *httpsFlag || os.Getenv("HTTPS") == "true"
	port := resolveStrFlag(portFlag, "port", "PORT")
	cert := resolveStrFlag(certFlag, "cert", "TLS_CERT")
	key := resolveStrFlag(keyFlag, "key", "TLS_KEY")

	if port == "" {
		if useHTTPS {
			port = "8443"
		} else {
			port = "8080"
		}
	}

	if useHTTPS {
		if err := router.RunTLS(":"+port, cert, key); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := router.Run(":" + port); err != nil {
			log.Fatal(err)
		}
	}
}
