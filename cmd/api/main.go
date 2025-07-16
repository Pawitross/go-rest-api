package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"pawrest/internal/api/middleware"
	"pawrest/internal/api/routes"
	"pawrest/internal/db"
	"pawrest/internal/yamlconfig"
)

type serverFlags struct {
	https *bool
	port  *string
	cert  *string
	key   *string
}

func main() {
	httpsFlag := flag.Bool("https", false, "Start the server with HTTPS")
	portFlag := flag.String("port", "", "Server port")
	certFlag := flag.String("cert", "keys/server.pem", "TLS certificate file location")
	keyFlag := flag.String("key", "keys/server.key", "TLS private key file location")
	flag.Parse()

	flags := serverFlags{
		https: httpsFlag,
		port:  portFlag,
		cert:  certFlag,
		key:   keyFlag,
	}

	if err := run(flags); err != nil {
		log.Fatal(err)
	}
}

func run(flags serverFlags) error {
	log.Println("Parsing env.yaml file...")
	cfg, err := yamlconfig.Parse("env.yaml")
	if err != nil {
		return err
	}

	log.Println("Connecting to the database...")
	database, err := db.ConnectToDB(cfg)
	if err != nil {
		return err
	}
	defer database.CloseDB()

	if err := middleware.InitLogger(); err != nil {
		return fmt.Errorf("failed to initialize logging middleware: %v", err)
	}
	defer middleware.CloseLogger()

	log.Println("Starting up the server...")
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(middleware.FileLogger())
	routes.Router(router, database, cfg)

	useHTTPS := *flags.https || os.Getenv("HTTPS") == "true"
	port := resolveStrFlag(flags.port, "port", "PORT")
	cert := resolveStrFlag(flags.cert, "cert", "TLS_CERT")
	key := resolveStrFlag(flags.key, "key", "TLS_KEY")

	if port == "" {
		if useHTTPS {
			port = "8443"
		} else {
			port = "8080"
		}
	}

	log.Printf("Started listening on port %v...\n", port)
	if useHTTPS {
		if err := router.RunTLS(":"+port, cert, key); err != nil {
			return err
		}
	} else {
		if err := router.Run(":" + port); err != nil {
			return err
		}
	}

	return nil
}

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
