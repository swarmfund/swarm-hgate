package hgate

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/hgate/server"
	"gitlab.com/swarmfund/hgate/server/problem"
	"gitlab.com/swarmfund/hgate/tx"
	"gitlab.com/swarmfund/horizon-connector"
)

type PaymentAction struct {
	server.Handler
	sub           *tx.Submitter
	sourceAccount string
	requestData   *tx.PaymentRequest
	paymentOp     *xdr.PaymentOp
}

func InitPaymentHandler(app *App) func(r *http.Request) server.HandlerI {
	return func(r *http.Request) server.HandlerI {
		app.Log.WithField("path", r.URL.Path).Info("Started request")
		return &PaymentAction{
			Handler: server.Handler{
				R:           r,
				Log:         app.Log.WithField("action", "payment"),
				Method:      http.MethodPost,
				ContentType: server.MimeJSON,
			},

			sub:           app.Sub,
			sourceAccount: app.Conf.KP.Address(),
		}
	}
}

func (action *PaymentAction) Post(w http.ResponseWriter, r *http.Request) {
	action.W = w
	action.Do(
		action.ValidateRequest,
		action.LoadRequest,
		action.GetParticipantsBalanceIDs,
		action.GetParticipantsFees,
		action.SubmitPaymentTx,
	)

	action.Render(w, problem.Success)
}

func (action *PaymentAction) LoadRequest() {
	action.requestData = new(tx.PaymentRequest)
	err := action.GetJSON(action.requestData)
	if err != nil {
		action.SetInvalidField("body", err)
		return
	}

	ok := IsValidAccountId(action.requestData.DestinationID)
	if !ok {
		action.SetInvalidField("destination_id", errors.New("must be valid accountId"))
		return
	}

	action.paymentOp = new(xdr.PaymentOp)
	action.paymentOp.Subject, err = ToXDRString256(action.requestData.Subject)
	if err != nil {
		action.SetInvalidField("subject", err)
		return
	}

	action.paymentOp.Reference, err = ToXDRString64(action.requestData.Reference)
	if err != nil {
		action.SetInvalidField("reference", err)
		return
	}

	tmpIntV, err := amount.Parse(action.requestData.Amount)
	if err != nil {
		action.SetInvalidField("amount", errors.New("must be a string that represents a number with four decimal places"))
		return
	}
	action.paymentOp.Amount = xdr.Int64(tmpIntV)
}

func (action *PaymentAction) GetParticipantsFees() {
	sourceFee, err := action.sub.GetPaymentFeeData(
		action.sourceAccount,
		action.requestData.Amount,
		action.requestData.Asset)
	if err != nil {
		action.Log.WithError(err).Error("unable to load fee for source")
		action.Err = err
		return
	}

	destinationFee, err := action.sub.GetPaymentFeeData(
		action.requestData.DestinationID,
		action.requestData.Amount,
		action.requestData.Asset)
	if err != nil {
		action.Log.WithError(err).Error("unable to load fee for destination")
		action.Err = err
		return
	}

	action.paymentOp.FeeData = xdr.PaymentFeeData{
		SourceFee:         *sourceFee,
		DestinationFee:    *destinationFee,
		SourcePaysForDest: action.requestData.PayFeeInsteadDest,
	}
}

func (action *PaymentAction) GetParticipantsBalanceIDs() {
	var err error

	action.paymentOp.SourceBalanceId, err = action.sub.GetBalanceIDForAsset(
		action.requestData.Asset,
		action.sourceAccount,
	)
	if err != nil {
		action.SetInvalidField("source_account", err)
		return
	}

	action.paymentOp.DestinationBalanceId, err = action.sub.GetBalanceIDForAsset(
		action.requestData.Asset,
		action.requestData.DestinationID,
	)
	if err != nil {
		action.SetInvalidField("destination_id", err)
		return
	}
}

func (action *PaymentAction) SubmitPaymentTx() {
	txSuccess, err := action.sub.CreatePaymentTx(action.paymentOp)
	if txSuccess != nil {
		// txSuccess is not nil only when tx
		// successfully submitted with 200 result code
		return
	}

	serr, ok := err.(horizon.SubmitError)
	if !ok {
		action.Log.WithError(err).Error("unable to submit tx")
		action.Err = &problem.ServerError
		return
	}

	action.RenderRawJSON(action.W, serr.ResponseBody())

}
