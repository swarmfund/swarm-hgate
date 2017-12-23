package tx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

type Balance struct {
	BalanceID string `json:"balance_id"`
	AccountID string `json:"account_id"`
	Asset     string `json:"asset"`
}

type BalanceResponse []Balance

func (b *Balance) XDRBalanceId() (xb xdr.BalanceId, err error) {
	xb, err = horizon.ParseBalanceID(b.BalanceID)
	if err != nil {
		err = fmt.Errorf("failed to convert balance_id into xdr: %s", err.Error())
	}
	return
}

func (b BalanceResponse) GetForAsset(asset string) (xb xdr.BalanceId, err error) {
	for _, bal := range b {
		if bal.Asset == asset {
			return bal.XDRBalanceId()
		}
	}
	err = errors.New(fmt.Sprintf("balanceID for asset %s is not present", asset))
	return
}

func (sub *Submitter) GetBalanceIDForAsset(asset string, accountId string) (xdr.BalanceId, error) {
	br, err := sub.GetBalances(accountId)
	if err != nil {
		return xdr.BalanceId{}, err
	}
	if br == nil {
		err = fmt.Errorf("account %s has no balances", accountId)
		return xdr.BalanceId{}, err
	}

	return br.GetForAsset(asset)

}

func (sub *Submitter) GetBalances(accountId string) (*BalanceResponse, error) {
	balancesPath := fmt.Sprintf("/accounts/%s/balances", accountId)

	req, err := horizon.NewRequest(sub.HorizonUrl, http.MethodGet, balancesPath)
	if err != nil {
		return nil, err
	}

	balances := new(BalanceResponse)
	err = sub.doRequest(req, balances)
	if err != nil && strings.Contains(err.Error(), "404") {
		return nil, errors.New("account not found")
	}
	return balances, err
}
