package payments

import "sync"

type inMemoryPaymentStore struct {
	mutex    *sync.Mutex
	payments []*Payment
}

func NewInMemoryPaymentStore() *inMemoryPaymentStore {
	return &inMemoryPaymentStore{
		payments: []*Payment{},
		mutex:    &sync.Mutex{},
	}
}

func (store *inMemoryPaymentStore) Create(payment *Payment) (string, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	err := ensurePaymentHasId(payment)
	if err != nil {
		return "", err
	}
	store.payments = append(store.payments, payment)
	return payment.Id, nil
}

func (store *inMemoryPaymentStore) Update(payment *Payment) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	err := store.Delete(payment.Id)
	if err != nil {
		return err
	}
	_, err = store.Create(payment)
	return err
}

func (store *inMemoryPaymentStore) Delete(id string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for i, payment := range store.payments {
		if payment.Id == id {
			store.payments = append(store.payments[:i], store.payments[i+1:]...)
			return nil
		}
	}
	return NoSuchElementErr
}

func (store *inMemoryPaymentStore) Count() (int, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return len(store.payments), nil
}

func (store *inMemoryPaymentStore) List(offset int, pageSize int) ([]*Payment, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	startIndex := offset * pageSize
	endIndex := startIndex + pageSize
	if endIndex > len(store.payments) {
		endIndex = len(store.payments)
	}
	return store.payments[startIndex:endIndex], nil
}

func (store *inMemoryPaymentStore) FindById(id string) (*Payment, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for _, payment := range store.payments {
		if payment.Id == id {
			return payment, nil
		}
	}

	return nil, NoSuchElementErr
}
