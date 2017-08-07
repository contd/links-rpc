package main_test

import (
  "log"
	"net"
  "os"
  "io"

  "golang.org/x/net/context"
  "google.golang.org/grpc"
	"github.com/contd/links-rpc/server"
  pb "github.com/contd/links-rpc/links"
  "testing"
)

const (
	port = ":5051"
  dbpath = "saved_test.sqlite"
)

func Server() {
  lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
  pb.RegisterLinksServer(s, main.NewServer(dbpath))
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}

func TestMain(m *testing.M) {
  go Server()
  os.Exit(m.Run())
}

func TestMessages(t *testing.T) {
    // Set up a connection to the Server.
    const address = "localhost:5051"
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
      t.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    // Creates a new LinkClient
  	client := pb.NewLinksClient(conn)

    // Test GetLink
    t.Run("GetLink", func(t *testing.T) {
        res, err := client.GetLink(context.Background(), &pb.LinkRequest{ Id: 65 })
        if err != nil {
          t.Fatalf("could not get link: %v", err)
        }
        t.Logf("Link: %v", res)
        if res.Id != 65 {
          t.Error("Expected 65, got ", res.Id)
        }
        if res.Url != "https://www.awsadvent.com/2016/12/06/just-add-code-fun-with-terraform-modules-and-aws/" {
          t.Error("Expected 'https://www.awsadvent.com/2016/12/06/just-add-code-fun-with-terraform-modules-and-aws/' got ", res.Url)
        }
        if res.Category != "tutorial" {
          t.Error("Expected 'tutorial', got ", res.Category)
        }
        if res.Created != "2017-06-11 00:36:13.696" {
          t.Error("Expected '2017-06-11 00:36:13.696', got ", res.Created)
        }
        if res.Done != 0 {
          t.Error("Expected 0, got ", res.Done)
        }

    })

    // Test GetLinks
    t.Run("GetLinks", func(t *testing.T) {
      filter := &pb.LinksFilter{Keyword: ""}
      stream, err := client.GetLinks(context.Background(), filter)
      if err != nil {
    		t.Fatalf("Error on get links: %v", err)
    	}
      for {
    		link, err := stream.Recv()
    		if err == io.EOF {
    			break
    		}
    		if err != nil {
    			t.Fatalf("%v.GetLinks(_) = _, %v", client, err)
    		}
    		t.Logf("Link: %v", link)
    	}
    })

    // Test CreateLink
    t.Run("CreateLink", func(t *testing.T) {
      link := &pb.LinkRequest{
    		Url:  "http://somerandome.com/url/path",
    		Category: "Javascript",
    		Created: "2017-03-30 13:36:01",
    		Done: 1,
    	}
      res, err := client.CreateLink(context.Background(), link)
    	if err != nil {
    		t.Fatalf("Could not create Link: %v", err)
    	}
      t.Logf("Link added: %d", res.Id)
    	if !res.Success {
    		t.Error("Expected true, got ", res.Success)
    	}
    })

    // Test UpdateLink
    t.Run("UpdateLink", func(t *testing.T) {
      link := &pb.LinkRequest{
    		Id: 124,
    		Url:  "http://somerandome.com/url/path2",
    		Category: "Javascript2",
    		Created: "2017-02-30 13:36:01",
    		Done: 0,
    	}
      res, err := client.UpdateLink(context.Background(), link)
    	if err != nil {
    		t.Fatalf("Could not update Link: %v", err)
    	}
      t.Logf("Link updated: %d", res.Id)
    	if !res.Success {
    		t.Error("Expected true, got ", res.Success)
    	}
    })

    // Test DeleteLink
    t.Run("DeleteLink", func(t *testing.T) {
      link := &pb.LinkRequest{
    		Id: 125,
    	}
      res, err := client.DeleteLink(context.Background(), link)
    	if err != nil {
    		t.Fatalf("Could not delete Link: %v", err)
    	}
      t.Logf("Link deleted: %d", res.Id)
    	if !res.Success {
    		t.Error("Expected true, got ", res.Success)
    	}
    })
}
