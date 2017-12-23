package main

import (
	"errors"

	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/xdr"
)

func isValidAccountId(accountId string) bool {
	_, err := keypair.Parse(accountId)
	return err == nil
}

func toXDRString64(s string) (xdr.String64, error) {
	if len(s) > 64 {
		return "", errors.New("invalid length - must be lower then 64")
	}
	return xdr.String64(s), nil
}
func toXDRString256(s string) (xdr.String256, error) {
	if len(s) > 256 {
		return "", errors.New("invalid length - must be lower then 256")
	}
	return xdr.String256(s), nil
}
