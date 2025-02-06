package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	pb "github.com/jasmine-nguyen/go-microservices/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	type user struct {
		userID   string
		password string
	}

	var u user

	statement, err := impl.db.Prepare("SELECT user_id, password FROM user WHERE user_id=? AND password=?")
	if err != nil {
		log.Println(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = statement.QueryRow(credentials.GetUserName(), credentials.GetPassword()).Scan(&u.userID, &u.password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	statement.Close()

	jwToken, err := createJWT(u.userID)
	if err != nil {
		return nil, err
	}

	return &pb.Token{Jwt: jwToken}, nil
}

func createJWT(userID string) (string, error) {
	key := []byte(os.Getenv("SIGNING_KEY"))
	log.Println("---signing key when creating jwt: ", key)
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-service",
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
	})

	signedToken, err := token.SignedString(key)
	log.Println("---signed token: ", signedToken)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return signedToken, nil
}

func (impl *Implementation) ValidateToken(ctx context.Context, token *pb.Token) (*pb.User, error) {
	key := []byte(os.Getenv("SIGNING_KEY"))
	log.Println("---signing key when validating jwt: ", key)
	log.Println("---jwt: ", token.Jwt)

	userId, err := validateJWT(token.Jwt, key)
	log.Println("---userId after validating jwt: ", userId)
	if err != nil {
		log.Println("---validate jwt error: ", err.Error())
		return nil, err
	}

	return &pb.User{UserId: userId}, nil
}

func validateJWT(t string, signingKey []byte) (string, error) {
	type MyClaims struct {
		jwt.RegisteredClaims
	}

	parsedToken, err := jwt.ParseWithClaims(t, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		log.Println("---error parsing token: ", err.Error())
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", status.Error(codes.Unauthenticated, "token expired")
		}

		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}

	log.Println("---parsedToken: ", parsedToken)
	claims, ok := parsedToken.Claims.(*MyClaims)
	if !ok {
		return "", status.Error(codes.Internal, "claims type assertion failed")
	}

	log.Println("---claims: ", claims)
	return claims.RegisteredClaims.Subject, nil
}
