package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

const (
	defaultServicePort  = 55555
	defaultServiceAddr  = ""
	defaultPageSize     = 20
	initBasketCapacity  = 200
	maxBasketCapacity   = 2000
	defaultDatabaseType = DbTypeMemory
	serviceOldAPIPath   = "baskets"
	serviceAPIPath      = "api"
	serviceUIPath       = "web"
	serviceName         = "request-baskets"
	basketNamePattern   = `^[\w\d\-_\.]{1,250}$`
	sourceCodeURL       = "https://github.com/milkinsa/QHRbaskets"
)

// ServerConfig describes server configuration.
type ServerConfig struct {
	ServerPort   int
	ServerAddr   string
	InitCapacity int
	MaxCapacity  int
	PageSize     int
	MasterToken  string
	DbType       string
	DbFile       string
	DbConnection string
	Baskets      []string
}

type arrayFlags []string

func (v *arrayFlags) String() string {
	return strings.Join(*v, ",")
}

func (v *arrayFlags) Set(value string) error {
	*v = append(*v, value)
	return nil
}

 func GetPort() int {
	 var port = os.Getenv("PORT")
     // Set a default port if there is nothing in the environment
     if port == "" {
         port = "4747"
         fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
     }
	 tempport, err :=strconv.Atoi(port)
	 if err != nil {
        // handle error
        fmt.Println(err)
    }
     return tempport
 }

 func GetDtbUrl() string {
	 
	dns := os.Getenv("DATABASE_URL")
	if dns == "" {
        log.Fatal("$DATABASE_URL is not set")
	   }
     return dns
 }

// CreateConfig creates server configuration base on application command line arguments
func CreateConfig() *ServerConfig {
	var port = flag.Int("p", GetPort(), "HTTP service port")
	var address = flag.String("l", defaultServiceAddr, "HTTP listen address")
	var initCapacity = flag.Int("size", initBasketCapacity, "Initial basket size (capacity)")
	var maxCapacity = flag.Int("maxsize", maxBasketCapacity, "Maximum allowed basket size (max capacity)")
	var pageSize = flag.Int("page", defaultPageSize, "Default page size")
	var masterToken = flag.String("token", "", "Master token, random token is generated if not provided")
	var dbType = flag.String("db", defaultDatabaseType, fmt.Sprintf(
		"Baskets storage type: %s - in-memory, %s - Bolt DB, %s - SQL database", DbTypeMemory, DbTypeBolt, DbTypeSQL))
	var dbFile = flag.String("file", "./baskets.db", "Database location, only applicable for file or SQL databases")
	var dbConnection = flag.String("conn", GetDtbUrl(), "Database connection string for SQL databases, if undefined \"file\" argument is considered")

	var baskets arrayFlags
	flag.Var(&baskets, "basket", "Name of a basket to auto-create during service startup (can be specified multiple times)")

	flag.Parse()

	var token = *masterToken
	if len(token) == 0 {
		token, _ = GenerateToken()
		log.Printf("[info] generated master token: %s", token)
	}

	return &ServerConfig{
		ServerPort:   *port,
		ServerAddr:   *address,
		InitCapacity: *initCapacity,
		MaxCapacity:  *maxCapacity,
		PageSize:     *pageSize,
		MasterToken:  token,
		DbType:       *dbType,
		DbFile:       *dbFile,
		DbConnection: *dbConnection,
		Baskets:      baskets}
}
