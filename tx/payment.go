package tx

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

type PaymentRequest struct {
	DestinationID     string `json:"destination_id"`
	Amount            string `json:"amount"`
	Asset             string `json:"asset"`
	Subject           string `json:"subject"`
	Reference         string `json:"reference"`
	PayFeeInsteadDest bool   `json:"pay_fee_instead_dest"`
}

func (sub *Submitter) CreatePaymentTx(paymentOp *xdr.PaymentOp) (*horizon.TransactionSuccess, error) {
	var accountID xdr.AccountId
	err := accountID.SetAddress(sub.KP.Address())
	if err != nil {
		return nil, errors.Wrap(err, "unable to set address")
	}

	op := xdr.Operation{
		SourceAccount: &accountID,
		Body: xdr.OperationBody{
			Type:      xdr.OperationTypePayment,
			PaymentOp: paymentOp,
		},
	}

	return sub.SubmitTx(op)
}
