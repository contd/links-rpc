package main

import (
	"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/contd/links-rpc/links"
)

const (
	address = "localhost:5051"
)

// getLink calls the RPC method GetLink of LinksServer
func getLink(client pb.LinksClient, link *pb.LinkRequest) {
	resp, err := client.GetLink(context.Background(), link)
	if err != nil {
		log.Fatalf("Could not get Link: %v", err)
	}
	if resp.Success {
		log.Printf("Link: %v", resp)
	}
}

// getLinks calls the RPC method GetLinks of LinksServer
func getLinks(client pb.LinksClient, filter *pb.LinksFilter) {
	// calling the streaming API
	stream, err := client.GetLinks(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get links: %v", err)
	}
	for {
		link, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetLinks(_) = _, %v", client, err)
		}
		log.Printf("Link: %v", link)
	}
}

// createLink calls the RPC method CreateLink of LinksServer
func createLink(client pb.LinksClient, link *pb.LinkRequest) {
	resp, err := client.CreateLink(context.Background(), link)
	if err != nil {
		log.Fatalf("Could not create Link: %v", err)
	}
	if resp.Success {
		log.Printf("A new Link has been added with id: %d", resp.Id)
	}
}

// updateLink calls the RPC method CreateLink of LinksServer
func updateLink(client pb.LinksClient, link *pb.LinkRequest) {
	resp, err := client.UpdateLink(context.Background(), link)
	if err != nil {
		log.Fatalf("Could not update Link: %v", err)
	}
	if resp.Success {
		log.Printf("Link has been updated: %v", resp)
	}
}

// deleteLink calls the RPC method CreateLink of LinksServer
func deleteLink(client pb.LinksClient, link *pb.LinkRequest) {
	resp, err := client.DeleteLink(context.Background(), link)
	if err != nil {
		log.Fatalf("Could not delete Link: %v", err)
	}
	if resp.Success {
		log.Printf("Link has been deleted: %d ", resp.Id)
	}
}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Creates a new LinkClient
	client := pb.NewLinksClient(conn)

	// Link req for get Link
  link = &pb.LinkRequest{
    Id: 65,
  }
  // Get link
  getLink(client, link)
	// Filter with an empty Keyword
	filter := &pb.LinksFilter{Keyword: ""}
	getLinks(client, filter)

	// Link req to create
	link := &pb.LinkRequest{
		Url:  "http://somerandome.com/url/path",
		Category: "Javascript",
		Created: "2017-03-30 13:36:01",
		Done: 1,
	}
	// Create a new link
	createLink(client, link)

	// Another Link req to create
	link = &pb.LinkRequest{
		Url:  "http://another.com/some/path/2",
		Category: "Politics",
		Created: "2017-01-10 16:26:21",
		Done: 1,
	}
	// Create a new link
	createLink(client, link)

	// UpdateLink
	link = &pb.LinkRequest{
		Id: 124,
		Url:  "http://somerandome.com/url/path2",
		Category: "Javascript2",
		Created: "2017-02-30 13:36:01",
		Done: 0,
	}
	updateLink(client, link)

	// DeleteLink
	link := &pb.LinkRequest{
		Id: 125
	}
	deleteLink(client, link)
}
