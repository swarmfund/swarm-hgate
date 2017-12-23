package tx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
)

// Submitter is
type Submitter struct {
	connector  *horizon.Connector
	AccountID  string
	KP         *keypair.Full
	HorizonUrl string
}

func NewSubmitter(horizonUrl string, accountID string, kp *keypair.Full) (s *Submitter, err error) {
	s = &Submitter{
		AccountID:  accountID,
		KP:         kp,
		HorizonUrl: horizonUrl,
	}
	s.connector, err = horizon.NewConnector(horizonUrl)
	return s, err
}

func (sub *Submitter) SubmitTx(ops ...xdr.Operation) (*horizon.TransactionSuccess, error) {
	tb := sub.connector.Transaction(&horizon.TransactionBuilder{
		Source:     sub.KP,
		Operations: []xdr.Operation(ops),
	})

	env, err := tb.Sign(sub.KP).Marshal64()
	if err != nil {
		return nil, err
	}

	return sub.connector.SubmitTXVerbose(*env)
}

func (sub *Submitter) doRequest(req *http.Request, responseDest interface{}) error {
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failed with status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(responseDest)
}
