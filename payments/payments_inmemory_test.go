package payments

import (
	"fmt"
	"github.com/Pallinder/sillyname-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateInMemoryPayment(t *testing.T) {
	// Fixtures
	t.Parallel()
	store := NewInMemoryPaymentStore()

	// When
	_, err := store.Create(&Payment{})

	// Then
	assert.NoError(t, err)
	count, err := store.Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestFindInMemoryPaymentById(t *testing.T) {
	// Given
	t.Parallel()
	payment := Payment{
		PaymentType: sillyname.GenerateStupidName(),
	}
	store := NewInMemoryPaymentStore()
	id, err := store.Create(&payment)
	assert.NoError(t, err)

	// When
	fetchedPayment, err := store.FindById(id)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, payment.PaymentType, fetchedPayment.PaymentType)
}

func TestListEmptyInMemoryCollection(t *testing.T) {
	// Given
	t.Parallel()
	store := NewInMemoryPaymentStore()

	// When
	payments, err := store.List(0, 10)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, payments)
}

func TestListInMemoryPage(t *testing.T) {
	// Given
	t.Parallel()
	store := NewInMemoryPaymentStore()
	for i := 0; i < 6; i++ {
		_, err := store.Create(&Payment{PaymentType: fmt.Sprintf("%d", i)})
		assert.NoError(t, err)
	}

	// When
	payments, err := store.List(1, 2)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 2, len(payments))
	assert.Equal(t, "2", payments[0].PaymentType)
	assert.Equal(t, "3", payments[1].PaymentType)
}
