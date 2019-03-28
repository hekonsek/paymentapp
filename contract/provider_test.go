package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hekonsek/paymentapp/api"
	"github.com/hekonsek/paymentapp/payments"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/pacts", dir)

func TestProviderContract(t *testing.T) {
	pact := &dsl.Pact{
		Consumer: "PaymentAppConsumer",
		Provider: "PaymentAppProvider",
	}
	err := startServer()
	assert.NoError(t, err)

	// Fixtures
	payment := payments.Payment{
		Id:          "foo",
		PaymentType: "Payment",
	}
	paymentJson, err := json.Marshal(&payment)
	assert.NoError(t, err)
	_, err = http.Post("http://localhost:8080/payments", "application/json", bytes.NewBuffer(paymentJson))
	assert.NoError(t, err)

	_, err = pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: "http://localhost:8080",
		PactURLs:        []string{filepath.ToSlash(fmt.Sprintf("%s/paymentappconsumer-paymentappprovider.json", pactDir))},
		StateHandlers: types.StateHandlers{
			"payment with id foo exists": func() error {
				return nil
			},
		},
	})
	assert.NoError(t, err)
}

func startServer() error {
	a := api.ApiServer{
		Port:  8080,
		Store: payments.NewInMemoryPaymentStore(),
	}
	err := a.Start()
	if err != nil {
		return err
	}

	return nil
}
