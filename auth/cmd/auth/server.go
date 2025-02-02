package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasmine-nguyen/go-microservices/auth/internal/implementation/auth"
	pb "github.com/jasmine-nguyen/go-microservices/auth/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	dbDriver = "mysql"
	dbName   = "auth"
)

var db *sql.DB

func main() {
	var err error

	dbUser := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")

	// Open a database connection
	dsn := fmt.Sprintf("%s:%s@tcp(mysql-auth:3306)/%s", dbUser, dbPassword, dbName)

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
	pb.RegisterAuthServiceServer(grpcServer, authServerImplementation)

	// listen and serve
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen on port 9000: %v", err)
	}

	log.Printf("server listening at: %v", listener.Addr())
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("server failed to serve: %v", err)
	}
}
