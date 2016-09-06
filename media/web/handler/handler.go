package handler

import (
	/*"fmt"
	"log"*/

	crawl "github.com/kazoup/platform/crawler/srv/proto/crawler"
	//"github.com/micro/go-micro/client"

	//"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

var (
	CrawlClient crawl.CrawlClient
)

func CrawlerStatus(ws *websocket.Conn) {
	/*req := client.NewRequest("go.micro.srv.crawl", "Crawl.Status", &crawl.StatusRequest{})

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
	}*/

}
