package handler

import (
	"fmt"
	"log"

	crawl "github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/micro/go-micro/client"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

var (
	CrawlClient crawl.CrawlClient
)

func Status(ws *websocket.Conn) {
	log.Print("Recive message on websocket")
	//for {

	//	if err := websocket.JSON.Send(ws, "ping"); err != nil {
	//		log.Print(err)
	//		return
	//	}
	//	time.Sleep(2 * time.Second)
	//}
	////Recive message on websocket
	//if err := websocket.JSON.Receive(ws, &m); err != nil {
	//	fmt.Println("err", err)
	//	return
	//}
	//c := crawl.NewCrawlClient("go.micro.srv.desktop", nil)
	//
	// Create new request to service go.micro.srv.example, method Example.Call
	// Request can be empty as its actually ignored and merely used to call the handler
	req := client.NewRequest("go.micro.srv.desktop", "Crawl.Status", &crawl.StatusRequest{})

	stream, err := client.Stream(context.Background(), req)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if err := stream.Send(&crawl.StatusRequest{}); err != nil {
		fmt.Println("err:", err)
		return
	}
	for stream.Error() == nil {
		rsp := &crawl.StatusResponse{}
		err := stream.Recv(rsp)
		if err != nil {
			fmt.Println("recv err", err)
			break
		}
		if err := websocket.JSON.Send(ws, rsp); err != nil {
			log.Print(err)
			return
		}
		fmt.Println("Stream: rsp:", rsp.Counter)
	}

	if stream.Error() != nil {
		fmt.Println("stream err:", err)
		return
	}

	if err := stream.Close(); err != nil {
		fmt.Println("stream close err:", err)
	}

}
