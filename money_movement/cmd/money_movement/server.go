package main

import (
	"database/sql"
	"fmt"
	"log"

	mm "github.com/jasmine-nguyen/go-microservices/money_movement/internal/implementation"
	pb "github.com/jasmine-nguyen/go-microservices/money_movement/proto"
	"google.golang.org/grpc"
)

const (
	dbDriver   = "mysql"
	dbUser     = "root"
	dbPassword = "Admin123"
	dbName     = "money_movement"
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
	pb.RegisterMoneyMovementServiceServer(grpcServer, mm.NewMoneyMovementImplementation(db))
}
