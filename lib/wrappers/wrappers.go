package wrappers

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/cli"
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

	if name == "db" {
		service := micro.NewService(
			micro.Name(sn),
			micro.Version("latest"),
			micro.Metadata(md),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
			//micro.Client(NewKazoupClientWithXrayTrace(sess)),
			micro.Flags(
				cli.StringFlag{
					Name:   "elasticsearch_hosts",
					EnvVar: "ELASTICSEARCH_HOSTS",
					Usage:  "Comma separated list of elasticsearch hosts",
					Value:  "elasticsearch:9200",
				},
			),
			micro.Action(func(c *cli.Context) {
				//parts := strings.Split(c.String("elasticsearch_hosts"), ",")
				//elastic.Hosts = parts
			}),
		)

		service.Init(
			//micro.Client(NewKazoupClientWithXrayTrace(sess)),
			micro.WrapClient(
				ContextClientWrapper(service),
				/*awsxray.NewClientWrapper(opts...),*/
				KazoupClientWrap(),
			),
			micro.WrapSubscriber(
				ContextSubscriberWrapper(service),
				AuthSubscriberWrapper,
				NewAfterSubscriberWrapper(),
				NewQuotaSubscriberWrapper(sn),
				LogSubscriberWrapper(),
			),
			micro.WrapHandler(
				ContextHandlerWrapper(service),
				/*awsxray.NewHandlerWrapper(opts...), */
				AuthHandlerWrapper,
				NewAfterHandlerWrapper(),
				NewQuotaHandlerWrapper(sn),
				LogHandlerWrapper(),
			),
		)

		return service
	}

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
			KazoupClientWrap(),
		),
		micro.WrapSubscriber(
			ContextSubscriberWrapper(service),
			AuthSubscriberWrapper,
			NewAfterSubscriberWrapper(),
			NewQuotaSubscriberWrapper(sn),
			LogSubscriberWrapper(),
		),
		micro.WrapHandler(
			ContextHandlerWrapper(service),
			/*awsxray.NewHandlerWrapper(opts...), */
			AuthHandlerWrapper,
			NewAfterHandlerWrapper(),
			NewQuotaHandlerWrapper(sn),
			LogHandlerWrapper(),
		),
	)

	return service
}
