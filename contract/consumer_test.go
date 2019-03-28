package contract

import (
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestConsumerContract(t *testing.T) {
	pact := &dsl.Pact{
		Consumer: "PaymentAppConsumer",
		Provider: "PaymentAppProvider",
		Host:     "localhost",
	}
	defer pact.Teardown()

	var test = func() error {
		u := fmt.Sprintf("http://localhost:%d/payment/foo", pact.Server.Port)
		req, err := http.NewRequest("GET", u, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			return err
		}
		if _, err = http.DefaultClient.Do(req); err != nil {
			return err
		}

		u = fmt.Sprintf("http://localhost:%d/payment/foo", pact.Server.Port)
		req, err = http.NewRequest("DELETE", u, nil)
		assert.NoError(t, err)
		if _, err = http.DefaultClient.Do(req); err != nil {
			return err
		}

		return err
	}

	type Payment struct {
		Type string `json:"type" pact:"example=Payment"`
	}
	pact.
		AddInteraction().
		Given("payment with id foo exists").
		UponReceiving("a request to get payment with id foo").
		WithRequest(dsl.Request{
			Method:  "GET",
			Path:    dsl.String("/payment/foo"),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json; charset=utf-8")},
			Body:    dsl.Match(&Payment{}),
		})

	pact.
		AddInteraction().
		Given("payment with id foo exists").
		UponReceiving("a request to delete payment with id foo").
		WithRequest(dsl.Request{
			Method: "DELETE",
			Path:   dsl.String("/payment/foo"),
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json; charset=utf-8")},
		})

	err := pact.Verify(test)
	assert.NoError(t, err)
}
