package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-toschool/syracuse"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/go-toschool/syracuse/postgres"
	"github.com/jmoiron/sqlx"
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
	cc, err := cs.Citizens.Select()
	if err != nil {
		return nil, err
	}

	data := make([]*citizens.Citizen, 0)
	for _, c := range cc {
		data = append(data, &citizens.Citizen{
			Id:        c.ID,
			Email:     c.Email,
			Fullname:  c.Fullname,
			CreatedAt: c.CreatedAt.Unix(),
			UpdatedAt: c.UpdatedAt.Unix(),
		})
	}

	return &citizens.SelectResponse{
		Data: data,
	}, nil
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
	u, err := cs.Citizens.Get(gr.UserId)
	if err != nil {
		return nil, err
	}

	u.Email = gr.Data.Email
	u.Fullname = gr.Data.Fullname
	if err := cs.Citizens.Update(u); err != nil {
		return nil, err
	}

	return &citizens.UpdateResponse{
		Data: &citizens.Citizen{
			Id:        u.ID,
			Email:     u.Email,
			Fullname:  u.Fullname,
			CreatedAt: u.CreatedAt.Unix(),
			UpdatedAt: u.UpdatedAt.Unix(),
		},
	}, nil
}

// Delete ...
func (cs *CitizensService) Delete(ctx context.Context, gr *citizens.DeleteRequest) (*citizens.DeleteResponse, error) {
	u, err := cs.Citizens.Get(gr.UserId)
	if err != nil {
		return nil, err
	}

	if err := cs.Citizens.Delete(u); err != nil {
		return nil, err
	}

	return &citizens.DeleteResponse{
		Data: &citizens.Citizen{
			Id:        u.ID,
			Email:     u.Email,
			Fullname:  u.Fullname,
			CreatedAt: u.CreatedAt.Unix(),
			UpdatedAt: u.UpdatedAt.Unix(),
		},
	}, nil
}
