package stripe

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/globals"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

func init() {
	stripe.Key = globals.STRIPE_SECRET_KEY
}

// GetCustomer returns email and profile as string
func GetCustomer(sID string) (string, string, error) {
	c, err := customer.Get(sID, nil)
	if err != nil {
		return "", "", err
	}

	b, err := json.Marshal(c)
	if err != nil {
		return "", "", err
	}

	return c.Email, string(b), nil
}

func DeleteCustomer(sID string) error {
	_, err := customer.Del(sID)
	if err != nil {
		return err
	}

	return nil
}
