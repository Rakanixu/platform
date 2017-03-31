package wrappers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/wrapper/trace/awsxray"
)

func NewKazoupClientWithXrayTrace(sess *session.Session) client.Client {
	opts := []awsxray.Option{
		// Used as segment name
		awsxray.WithName("com.kazoup.client"),
		// Specify X-Ray Daemon Address
		// awsxray.WithDaemon("localhost:2000"),
		// Or X-Ray Client
		awsxray.WithClient(xray.New(sess, &aws.Config{Region: aws.String("eu-west-1")})),
	}
	return client.NewClient(
		client.WrapCall(awsxray.NewCallWrapper(opts...)),
	)
}
