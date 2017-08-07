package main

import (
  "log"
	"net"
  "strings"

  "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/contd/links-rpc/links"
)

const (
	port = ":5051"
  dbpath = "saved.sqlite"
)

// server is used to implement links.LinksServer.
type server struct {
  DB     *sqlx.DB
}

// GetLink returns link given Id
func (s *server) GetLink(ctx context.Context, req *pb.LinkRequest) (*pb.LinkResponse, error) {
  l := link{ID: req.Id}
	if err := l.getLink(s.DB); err != nil {
    return nil, err
  }
  return &pb.LinkResponse{Id: l.ID, Url: l.Url, Category: l.Category, Created: l.Created, Done: l.Done, Success: true}, nil
}

// GetLinks returns all links by given filter
func (s *server) GetLinks(filter *pb.LinksFilter, stream pb.Links_GetLinksServer) error {
  links, err := getLinks(s.DB)
	if err != nil {
    return err
  }
	for _, link := range links {
		if filter.Keyword != "" {
			if !strings.Contains(link.Category, filter.Keyword) {
				continue
			}
		}
		if err := stream.Send(&pb.LinkRequest{ Id: link.ID, Url: link.Url, Category: link.Category, Created: link.Created, Done: link.Done }); err != nil {
			return err
		}
	}
	return nil
}

// CreateLink creates a new Link
func (s *server) CreateLink(ctx context.Context, req *pb.LinkRequest) (*pb.LinkResponse, error) {
  l := link{
    Url: req.Url,
    Category: req.Category,
    Created: req.Created,
    Done: req.Done,
  }
  id, err := l.createLink(s.DB)
  if err != nil {
    return nil, err
  }
	return &pb.LinkResponse{Id: int32(id), Success: true}, nil
}

// UpdateLink updates a Link
func (s *server) UpdateLink(ctx context.Context, req *pb.LinkRequest) (*pb.LinkResponse, error) {
  l := link{
    Url: req.Url,
    Category: req.Category,
    Created: req.Created,
    Done: req.Done,
  }
  if err := l.updateLink(s.DB); err != nil {
    return nil, err
    //return &pb.LinkNotFoundResponse{Id: req.Id, Notfound: "Link ID not found!"}, nil
  }
	return &pb.LinkResponse{Id: req.Id, Success: true}, nil
}

// DeleteLink updates a Link
func (s *server) DeleteLink(ctx context.Context, req *pb.LinkRequest) (*pb.LinkResponse, error) {
  l := link{ID: req.Id}
	if err := l.deleteLink(s.DB); err != nil {
    return nil, err
    //return &pb.LinkNotFoundResponse{Id: req.Id, Notfound: "Link ID not found!"}, nil
  }
  return &pb.LinkResponse{Id: req.Id, Success: true}, nil
}

func (s *server) connDB(dbpath string) {
  var err error
  s.DB, err = sqlx.Connect("sqlite3", dbpath)
  if err != nil {
		log.Fatal(err)
	}
}

func NewServer(dbpath string) *server {
	s := new(server)
  // Open DB conn
  s.connDB(dbpath)
	return s
}

func main() {
  lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
  pb.RegisterLinksServer(s, NewServer(dbpath))
	s.Serve(lis)
}
