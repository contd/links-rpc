package mock_links

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	linksmock "github.com/contd/links-rpc/mock_links"
	linkspb "github.com/contd/links-rpc/links"
)

var (
	msg = &linkspb.LinkResponse{
		Id:   65,
		Url:  "https://www.awsadvent.com/2016/12/06/just-add-code-fun-with-terraform-modules-and-aws/",
    Category: "tutorial",
    Created: "2017-06-11 00:36:13.696",
    Done: 0,
    Success: true,
	}
)

func TestGetLinks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock for the stream returned by GetLinks
	stream := linksmock.NewMockLinks_GetLinksClient(ctrl)
	// set expectation on sending.
	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil)
	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(msg, nil)
	stream.EXPECT().CloseSend().Return(nil)
	// Create mock for the client interface.
	linksclient := linksmock.NewMockLinksClient(ctrl)
	// Set expectation on GetLinks
	linksclient.EXPECT().GetLinks(
		gomock.Any(),
	).Return(stream, nil)
	if err := testGetLinks(linksclient); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func testGetLinks(client linkspb.LinksClient) error {
	stream, err := client.GetLinks(context.Background())
	if err != nil {
		return err
	}
	if err := stream.Send(msg); err != nil {
		return err
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	got, err := stream.Recv()
	if err != nil {
		return err
	}
	if !proto.Equal(got, msg) {
		return fmt.Errorf("stream.Recv() = %v, want %v", got, msg)
	}
	return nil
}
