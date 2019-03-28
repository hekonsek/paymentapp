package payments

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"time"
)

// Payment model

type Payment struct {
	PaymentType    string     `json:"type"`
	Id             string     `json:"id"`
	Version        string     `json:"version"`
	OrganisationId string     `json:"organisation_id"`
	Attributes     Attributes `json:"attributes"`
}

// ProcessingDate is a wrapper around time.Time providing custom data formatting (YYYY-MM-DD) for
// payment.attributes.processing_date JSON path.
type ProcessingDate time.Time

func (t ProcessingDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

type Attributes struct {
	Amount               *big.Float         `json:"amount"`
	BeneficiaryParty     BeneficiaryParty   `json:"beneficiary_party"`
	ChargesInformation   ChargesInformation `json:"charges_information"`
	Currency             string             `json:"currency"`
	EndToEndReference    string             `json:"end_to_end_reference"`
	NumericReference     string             `json:"numeric_reference"`
	PaymentId            string             `json:"payment_id"`
	PaymentPurpose       string             `json:"payment_purpose"`
	PaymentScheme        string             `json:"payment_scheme"`
	PaymentType          string             `json:"payment_type"`
	ProcessingDate       ProcessingDate     `json:"processing_date"`
	Reference            string             `json:"reference"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type"`
	SchemePaymentType    string             `json:"scheme_payment_type"`
	DebtorParty          DebtorParty        `fx:"debtor_party"`
	Fx                   Fx                 `json:"fx"`
	SponsorParty         SponsorParty       `json:"sponsor_party"`
}

type DebtorParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	Address           string `json:"address"`
	BankId            string `json:"bank_id"`
	BankIdCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

type Fx struct {
	ContractReference string     `json:"contract_reference"`
	ExchangeRate      string     `json:"exchange_rate"`
	OriginalAmount    *big.Float `json:"original_amount"`
	OriginalCurrency  string     `json:"original_currency"`
}

type SponsorParty struct {
	AccountNumber string `json:"account_number"`
	BankId        string `json:"bank_id"`
	BankIdCode    string `json:"bank_id_code"`
}

type BeneficiaryParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankId            string `json:"bank_id"`
	BankIdCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

type ChargesInformation struct {
	BearerCode              string          `json:"bearer_code"`
	ReceiverChargesAmount   big.Float       `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string          `json:"receiver_charges_currency"`
	ChargesInformation      []SenderCharges `json:"charges_information"`
}

type SenderCharges struct {
	Amount   big.Float `json:"amount"`
	Currency string    `json:"currency"`
}

// Payment store

var NoSuchElementErr = errors.New("no element with given ID")

type PaymentStore interface {
	Create(payment *Payment) (id string, err error)
	Update(payment *Payment) (err error)
	Delete(id string) error
	Count() (int, error)
	List(offset int, pageSize int) ([]*Payment, error)
	FindById(id string) (*Payment, error)
}

func ensurePaymentHasId(payment *Payment) error {
	if payment.Id == "" {
		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		payment.Id = u.String()
	}

	return nil
}
