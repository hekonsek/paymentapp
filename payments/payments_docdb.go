package payments

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
	"strings"
	"time"
)

type docdbPaymentStore struct {
	mongo      *mongo.Client
	collection *mongo.Collection
}

func NewDocdbPaymentStore(host string, port int) (*docdbPaymentStore, error) {
	if host == "" {
		host = os.Getenv("AWSDOCDB_SERVICE_HOST")
		if host == "" {
			host = "localhost"
		}
	}
	if port < 1 {
		portFromEnv, err := strconv.Atoi(os.Getenv("AWSDOCDB_SERVICE_PORT"))
		if err != nil {
			port = 27017
		}
		port = portFromEnv
	}

	store := &docdbPaymentStore{}
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port)))
	if err != nil {
		return nil, err
	}
	store.mongo = client
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	store.collection = store.mongo.Database("paymentsapp").Collection("payments")

	return store, nil
}

func (store *docdbPaymentStore) Create(payment *Payment) (id string, err error) {
	err = ensurePaymentHasId(payment)
	if err != nil {
		return "", err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	paymentBson, err := bson.Marshal(payment)
	if err != nil {
		return "", err
	}
	_, err = store.collection.InsertOne(ctx, paymentBson)
	if err != nil {
		return "", err
	}

	return payment.Id, nil
}

func (store *docdbPaymentStore) Update(p *Payment) (err error) {
	paymentBson, err := bson.Marshal(p)
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	results := store.collection.FindOneAndReplace(ctx, bson.M{"id": p.Id}, paymentBson)
	if results.Err() != nil {
		return results.Err()
	}

	return nil
}

func (store *docdbPaymentStore) Delete(id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	results, err := store.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if results.DeletedCount > 1 {
		return errors.New("deleted too many documents")
	} else if results.DeletedCount == 0 {
		return NoSuchElementErr
	}
	return nil
}

func (store *docdbPaymentStore) Count() (int, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	count, err := store.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return -1, err
	}
	return int(count), nil
}

func (store *docdbPaymentStore) List(offset int, pageSize int) ([]*Payment, error) {
	startIndex := int64(offset) * int64(pageSize)
	limit := int64(pageSize)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := store.collection.Find(ctx, bson.M{}, &options.FindOptions{Skip: &startIndex, Limit: &limit})
	if err != nil {
		return nil, err
	}
	var payments []*Payment
	for cursor.Next(ctx) {
		var p *Payment
		err = cursor.Decode(&p)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (store *docdbPaymentStore) FindById(id string) (*Payment, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	results := store.collection.FindOne(ctx, bson.M{"id": id})
	if results.Err() != nil {
		return nil, results.Err()
	}

	var payment Payment
	err := results.Decode(&payment)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return nil, NoSuchElementErr
		}
		return nil, err
	}
	return &payment, nil
}
