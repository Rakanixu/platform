package main

import (
	"fmt"

	"github.com/kazoup/platform/ldap/srv/handler"
	"github.com/kazoup/platform/ldap/srv/proto"
	"github.com/micro/go-micro"
)

func main() {
	//Crete new service
	service := micro.NewService(
		micro.Name("go.micro.srv.ldap"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "template",
		}),
	)
	service.Init()
	ldap.RegisterLDAPHandler(service.Server(), new(handler.LDAP))
	//Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
