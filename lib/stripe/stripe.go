package stripe

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/globals"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
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

// SaveCard saves a credit card to a stripe user
func SaveCard(uID, tID string) error {
	cp := &stripe.CustomerParams{}
	if err := cp.SetSource(tID); err != nil {
		return err
	}

	_, err := customer.Update(uID, cp)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSubscription
func UpdateSubscription(sID, newPlan string) error {
	sp := &stripe.SubParams{
		Plan: newPlan,
	}

	_, err := sub.Update(sID, sp)
	if err != nil {
		return err
	}

	return nil
}
