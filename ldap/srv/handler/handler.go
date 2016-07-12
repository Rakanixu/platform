package handler

import (
	"log"

	client "github.com/jtblin/go-ldap-client"
	"github.com/kazoup/platform/ldap/srv/proto"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

//LDAP ...
type LDAP struct{}

//Bind ...
func (c *LDAP) Bind(
	ctx context.Context,
	req *ldap.BindRequest,
	res *ldap.BindResponse) error {
	//Check if req not empty
	log.Print(req)
	if req.Config == nil {
		return errors.BadRequest("com.kazoup.srv.ldap.Config.Create", "Configuration can not be blank")
	}
	//Create LDAP client
	client := &client.LDAPClient{
		Base:         req.Config.BaseDn,
		Host:         req.Config.Host,
		Port:         int(req.Config.Port),
		UseSSL:       req.Config.UseSsl,
		BindDN:       req.Config.BindDn,
		BindPassword: req.Config.BindDn,
		UserFilter:   req.Config.UserFilter,  // "(uid=%s)"
		GroupFilter:  req.Config.GroupFilter, //"(memberUid=%s)",
		Attributes:   req.Config.Attr,        //[]string{"givenName", "sn", "mail", "uid"},
	}
	//defer connection close
	defer client.Close()
	//Validate binding with AD
	err := client.Connect()
	if err != nil {
		return errors.BadRequest("com.kazoup.srv.ldap.Config.Create Can not bind to AD", err.Error())
	}
	return nil
}

//Login ...
func (c *LDAP) Login(
	ctx context.Context,
	req *ldap.LoginRequest,
	res *ldap.LoginResponse) error {
	//TODO implement
	return nil
}

//Search ...
func (c *LDAP) Search(
	ctx context.Context,
	req *ldap.SearchRequest,
	res *ldap.SearchResponse) error {
	//TODO implement
	return nil
}
