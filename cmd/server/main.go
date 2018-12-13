package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"github.com/go-toschool/syracuse"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/go-toschool/syracuse/postgres"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

func main() {
	port := flag.Int64("port", 8001, "listening port")
	postgresDSN := flag.String("postgres-dsn", "postgres://localhost:5432/syracuse?sslmode=disable", "Postgres DSN")

	flag.Parse()

	db, err := sqlx.Connect("postgres", *postgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	srv := grpc.NewServer()

	citizens.RegisterCitizenshipServer(srv, &CitizensService{
		Citizens: &postgres.CitizensService{
			Store: db,
		},
	})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting Syracuse service...")
	log.Println(fmt.Sprintf("Syracuse service, Listening on: %d", *port))
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

// CitizensService ...
type CitizensService struct {
	Citizens syracuse.Citizens
}

// Get ...
func (cs *CitizensService) Get(ctx context.Context, gr *citizens.GetRequest) (*citizens.GetResponse, error) {
	c, err := cs.Citizens.Get(gr.GetUserId())
	if err != nil {
		return nil, err
	}

	return &citizens.GetResponse{
		Data: &citizens.Citizen{
			Id:        c.ID,
			Email:     c.Email,
			Fullname:  c.Fullname,
			CreatedAt: c.CreatedAt.Unix(),
			UpdatedAt: c.UpdatedAt.Unix(),
		},
	}, nil
}

// Select ...
func (cs *CitizensService) Select(ctx context.Context, gr *citizens.SelectRequest) (*citizens.SelectResponse, error) {
	return nil, nil
}

// Create ...
func (cs *CitizensService) Create(ctx context.Context, gr *citizens.CreateRequest) (*citizens.CreateResponse, error) {
	c := &syracuse.Citizen{
		Email:    gr.Data.Email,
		Fullname: gr.Data.Fullname,
	}

	if err := cs.Citizens.Create(c); err != nil {
		return nil, err
	}

	return &citizens.CreateResponse{
		Data: &citizens.Citizen{
			Id:        c.ID,
			Email:     c.Email,
			Fullname:  c.Fullname,
			CreatedAt: c.CreatedAt.Unix(),
			UpdatedAt: c.UpdatedAt.Unix(),
		},
	}, nil
}

// Update ...
func (cs *CitizensService) Update(ctx context.Context, gr *citizens.UpdateRequest) (*citizens.UpdateResponse, error) {
	return nil, nil
}

// Delete ...
func (cs *CitizensService) Delete(ctx context.Context, gr *citizens.DeleteRequest) (*citizens.DeleteResponse, error) {
	return nil, nil
}
