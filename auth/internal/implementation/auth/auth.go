package auth

import (
	"context"
	"database/sql"

	pb "github.com/jasmine-nguyen/go-microservices/auth/proto"
)

type Implementation struct {
	db *sql.DB
	pb.UnimplementedAuthServiceServer
}

func NewAuthImplementation(db *sql.DB) *Implementation {
	return &Implementation{
		db: db,
	}
}

func (impl *Implementation) GetToken(ctx context.Context, credentials *pb.Credentials) (*pb.Token, error) {
	return &pb.Token{}, nil
}

func (impl *Implementation) ValidateToken(ctx context.Context, token *pb.Token) (*pb.User, error) {
	return &pb.User{}, nil
}
