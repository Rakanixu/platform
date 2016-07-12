package publisher

import "github.com/micro/go-micro/broker"

// PublishData publishes to the broker
func PublishData(m map[string]string, b []byte) error {
	msg := &broker.Message{
		Header: m,
		Body:   []byte(b),
	}
	err := broker.Publish(m["topic"], msg)

	// TODO: handle or log errors
	if err != nil {
		return err
	}

	return nil
}
