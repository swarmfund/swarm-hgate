package tx

import (
	"fmt"
	"net/http"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

type Fee struct {
	PaymentFee string `json:"payment_fee"`
	FixedFee   string `json:"fixed_fee"`
}

type FeeData struct {
	DestinationFee    Fee  `json:"destination_fee"`
	SourceFee         Fee  `json:"source_fee"`
	SourcePaysForDest bool `json:"source_pays_for_dest"`
}

func (fe *FeeEntry) ToXDRFeeData() (*xdr.FeeData, error) {
	feeData := new(xdr.FeeData)
	tmpIntV, err := amount.Parse(fe.Percent)
	if err != nil {
		return nil, err
	}
	feeData.PaymentFee = xdr.Int64(tmpIntV)

	tmpIntV, err = amount.Parse(fe.Fixed)
	if err != nil {
		return nil, err
	}
	feeData.FixedFee = xdr.Int64(tmpIntV)

	return feeData, nil
}

type FeeEntry struct {
	Asset       string `json:"asset"`
	Fixed       string `json:"fixed"`
	Percent     string `json:"percent"`
	FeeType     int    `json:"fee_type"`
	Subtype     int64  `json:"subtype"`
	AccountID   string `json:"account_id"`
	AccountType int32  `json:"account_type"`
	LowerBound  string `json:"lower_bound"`
	UpperBound  string `json:"upper_bound"`
}

func (s *Submitter) GetPaymentFeeData(account, amount, asset string) (*xdr.FeeData, error) {
	return s.GetFeeData(xdr.FeeTypePaymentFee, map[string]string{
		"account": account,
		"amount":  amount,
		"asset":   asset,
	})
}

func (s *Submitter) GetFeeData(feeType xdr.FeeType, filters map[string]string) (*xdr.FeeData, error) {
	feePath := fmt.Sprintf("/fees/%d/?", feeType)

	for key, val := range filters {
		feePath += fmt.Sprintf("%s=%s&", key, val)
	}

	req, err := horizon.NewRequest(s.HorizonUrl, http.MethodGet, feePath)
	if err != nil {
		return nil, err
	}

	feeEntry := new(FeeEntry)
	err = s.doRequest(req, feeEntry)
	if err != nil {
		return nil, err
	}

	return feeEntry.ToXDRFeeData()
}
