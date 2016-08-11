package main

import "github.com/micro/go-micro"

func main() {

	service := micro.NewService(
		micro.Name("go.micro.srv.bleve"),
		micro.Version("latest"),
	)
}
