package mm

import (
	"context"
	"database/sql"

	pb "github.com/jasmine-nguyen/go-microservices/proto"
)

type Implementation struct {
	db *sql.DB
	pb.UnimplementedMoneyMovementServiceServer
}

func NewMoneyMovementImplementation(db *sql.DB) *Implementation {
	return &Implementation{}
}

func (impl *Implementation) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	return &pb.AuthorizeResponse{}, nil
}
