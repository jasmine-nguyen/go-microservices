package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jasmine-nguyen/go-microservices/auth/internal/implementation/auth"
	pb "github.com/jasmine-nguyen/go-microservices/proto"
	"google.golang.org/grpc"
)

const (
	dbDriver   = "mysql"
	dbUser     = "root"
	dbPassword = "Admin123"
	dbName     = "User"
)

var db *sql.DB

func main() {
	var err error

	// Open a database connection
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPassword, dbName)

	db, err = sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("Error closing db: %s", err)
		}
	}()

	// Ping the database to ensure the connection is valid
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// grpc server setup
	grpcServer := grpc.NewServer()
	authServerImplementation := auth.NewAuthImplementation(db)
	pb.RegisterAuthServiceServer(authServerImplementation, grpcServer)
}
