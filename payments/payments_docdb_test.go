package payments

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePayment(t *testing.T) {
	// Given
	t.Parallel()
	payment := Payment{
		PaymentType: "foo",
	}
	store, err := NewDocdbPaymentStore("", -1)
	assert.NoError(t, err)

	// When
	_, err = store.Create(&payment)

	// Then
	assert.NoError(t, err)
	count, err := store.Count()
	assert.NoError(t, err)
	assert.True(t, count > 0)
}

func TestFindPaymentById(t *testing.T) {
	// Given
	t.Parallel()
	payment := Payment{
		PaymentType: "foo",
	}
	store, err := NewDocdbPaymentStore("", -1)
	assert.NoError(t, err)
	id, err := store.Create(&payment)
	assert.NoError(t, err)

	// When
	fetchedPayment, err := store.FindById(id)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "foo", fetchedPayment.PaymentType)
}

func TestUpdatePayment(t *testing.T) {
	// Given
	t.Parallel()
	payment := Payment{
		PaymentType: "foo",
	}
	store, err := NewDocdbPaymentStore("", -1)
	assert.NoError(t, err)
	id, err := store.Create(&payment)
	assert.NoError(t, err)
	payment.PaymentType = "bar"

	// When
	err = store.Update(&payment)

	// Then
	assert.NoError(t, err)
	fetchedPayment, err := store.FindById(id)
	assert.NoError(t, err)
	assert.Equal(t, "bar", fetchedPayment.PaymentType)
}

func TestNotFindPaymentById(t *testing.T) {
	// Given
	t.Parallel()
	store, err := NewDocdbPaymentStore("", -1)
	assert.NoError(t, err)

	// When
	_, err = store.FindById("noSuchId")

	// Then
	assert.Equal(t, NoSuchElementErr, err)
}

func TestDeleteFromDocDb(t *testing.T) {
	// Given
	t.Parallel()
	store, err := NewDocdbPaymentStore("", -1)
	assert.NoError(t, err)
	id, err := store.Create(&Payment{})
	assert.NoError(t, err)

	// When
	err = store.Delete(id)

	// Then
	assert.NoError(t, err)
	_, err = store.FindById("noSuchId")
	assert.Equal(t, NoSuchElementErr, err)
}

func TestListPayments(t *testing.T) {
	// Given
	t.Parallel()
	payment := Payment{}
	store, err := NewDocdbPaymentStore("", -1)
	assert.NoError(t, err)
	_, err = store.Create(&payment)
	assert.NoError(t, err)

	// When
	payments, err := store.List(0, 10)

	// Then
	assert.NoError(t, err)
	assert.True(t, len(payments) > 0)
}
