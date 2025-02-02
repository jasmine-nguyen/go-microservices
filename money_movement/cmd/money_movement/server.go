package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	mm "github.com/jasmine-nguyen/go-microservices/money_movement/internal/implementation"
	pb "github.com/jasmine-nguyen/go-microservices/money_movement/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	dbDriver   = "mysql"
	dbUser     = "money_movement_user"
	dbPassword = "Auth123"
	dbName     = "money_movement"
)

var db *sql.DB

func main() {
	var err error

	// Open a database connection
	dsn := fmt.Sprintf("%s:%s@tcp(mysql-money-movement:3306)/%s", dbUser, dbPassword, dbName)

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

	// listen and serve
	listener, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalf("failed to listen on port 7000: %v", err)
	}

	log.Printf("server listening at: %v", listener.Addr())
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("server failed to serve: %v", err)
	}
}
