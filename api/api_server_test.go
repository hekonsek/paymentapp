package api

import (
	"bytes"
	"github.com/Pallinder/sillyname-go"
	"github.com/gin-gonic/gin/json"
	"github.com/hekonsek/paymentapp/payments"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"
)
import "github.com/tidwall/gjson"

func init() {
	a := ApiServer{
		Port:  8080,
		Store: payments.NewInMemoryPaymentStore(),
	}
	err := a.Start()
	if err != nil {
		panic(err)
	}
}

func TestGenerateNewIdForCreatedPayment(t *testing.T) {
	// Given
	t.Parallel()
	payment := payments.Payment{
		PaymentType: sillyname.GenerateStupidName(),
		Attributes: payments.Attributes{
			Amount: big.NewFloat(float64(5.32)),
		},
	}
	paymentJson, err := json.Marshal(&payment)
	assert.NoError(t, err)

	// When
	resp, err := http.Post("http://localhost:8080/payments", "application/json", bytes.NewBuffer(paymentJson))
	defer func() {
		err = resp.Body.Close()
		assert.NoError(t, err)
	}()
	createdResponse, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Then
	id := gjson.GetBytes(createdResponse, "id").Str
	assert.NotEmpty(t, id)
}

func TestGetPaymentById(t *testing.T) {
	// Given
	t.Parallel()
	payment := payments.Payment{
		PaymentType: sillyname.GenerateStupidName(),
		Attributes: payments.Attributes{
			Amount: big.NewFloat(float64(5.32)),
		},
	}
	paymentJson, err := json.Marshal(&payment)
	assert.NoError(t, err)
	resp, err := http.Post("http://localhost:8080/payments", "application/json", bytes.NewBuffer(paymentJson))
	defer func() {
		err = resp.Body.Close()
		assert.NoError(t, err)
	}()
	createdResponse, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	id := gjson.GetBytes(createdResponse, "id").Str
	assert.NotEmpty(t, id)

	// When
	resp, err = http.Get("http://localhost:8080/payment/" + id)
	assert.NoError(t, err)
	paymentJson, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		assert.NoError(t, err)
	}()

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, id, gjson.GetBytes(paymentJson, "id").Str)
}

func TestListPayments(t *testing.T) {
	// Given
	t.Parallel()

	// When
	resp, err := http.Get("http://localhost:8080/payments")
	paymentsJson, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		assert.NoError(t, err)
	}()

	// Then
	id := gjson.GetBytes(paymentsJson, "data")
	assert.True(t, len(id.Array()) > -1)
}

func TestCountPayments(t *testing.T) {
	// Given
	t.Parallel()

	// When
	resp, err := http.Get("http://localhost:8080/payments/count")
	paymentsJson, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		assert.NoError(t, err)
	}()

	// Then
	count := gjson.GetBytes(paymentsJson, "count").Int()
	assert.True(t, count >= 0)
}

func TestDelete(t *testing.T) {
	// Given
	t.Parallel()
	payment := payments.Payment{
		PaymentType: sillyname.GenerateStupidName(),
		Attributes: payments.Attributes{
			Amount: big.NewFloat(float64(5.32)),
		},
	}
	paymentJson, err := json.Marshal(&payment)
	assert.NoError(t, err)
	resp, err := http.Post("http://localhost:8080/payments", "application/json", bytes.NewBuffer(paymentJson))
	defer func() {
		err = resp.Body.Close()
		assert.NoError(t, err)
	}()
	createdResponse, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	id := gjson.GetBytes(createdResponse, "id").Str

	// When
	req, err := http.NewRequest("DELETE", "http://localhost:8080/payment/"+id, nil)
	assert.NoError(t, err)
	resp, err = (&http.Client{}).Do(req)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, 200, resp.StatusCode)
	resp, err = http.Get("http://localhost:8080/payment/" + id)
	assert.NoError(t, err)
	defer func() {
		err = resp.Body.Close()
		assert.NoError(t, err)
	}()
	assert.Equal(t, 404, resp.StatusCode)
}
