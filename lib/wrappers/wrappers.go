package wrappers

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro"
	"github.com/micro/go-os/monitor"
	"log"
	"os"
	"time"
)

func NewKazoupService(name string, mntr ...monitor.Monitor) micro.Service {
	var m monitor.Monitor

	// Check if monitor available
	if len(mntr) > 0 && mntr[0] != nil {
		m = mntr[0]
	}

	//Get AWS session credentials in ~/.aws/credentials
	/*	sess, err := session.NewSession()
		if err != nil {
			log.Fatal(err)
		}
		// New AWSXRAY with default
		//x := xray.New(sess, &aws.Config{Region: aws.String("eu-west-1")})
		opts := []awsxray.Option{
			// Used as segment name
			awsxray.WithName(name),
			// Specify X-Ray Daemon Address
			// awsxray.WithDaemon("localhost:2000"),
			// Or X-Ray Client
			awsxray.WithClient(xray.New(sess, &aws.Config{Region: aws.String("eu-west-1")})),
		}*/

	//FIXME hacked we should just pass micro.Flags
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("oops can't get hostname")
	}

	md := map[string]string{
		"hostname": hostname,
	}

	sn := fmt.Sprintf("%s.srv.%s", globals.NAMESPACE, name)

	var service micro.Service
	if m == nil {
		service = micro.NewService(
			micro.Name(sn),
			micro.Version("latest"),
			micro.Metadata(md),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
		)
	} else {
		service = micro.NewService(
			micro.Name(sn),
			micro.Version("latest"),
			micro.Metadata(md),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
		)
	}

	service.Init(
		//micro.Client(NewKazoupClientWithXrayTrace(sess)),
		micro.WrapClient(
			ContextClientWrapper(service),
			/*awsxray.NewClientWrapper(opts...),*/
		),
		micro.WrapSubscriber(
			NewContextSubscriberWrapper(service),
			NewAuthSubscriberWrapper(),
			NewAfterSubscriberWrapper(),
			NewQuotaSubscriberWrapper(sn),
			NewLogSubscriberWrapper(),
		),
		micro.WrapHandler(
			NewContextHandlerWrapper(service),
			/*awsxray.NewHandlerWrapper(opts...), */
			NewAuthHandlerWrapper(),
			NewAfterHandlerWrapper(),
			NewQuotaHandlerWrapper(sn),
			NewLogHandlerWrapper(),
		),
	)

	return service
}
