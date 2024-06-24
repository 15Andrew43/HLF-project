package handler

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

func GetClient(ctx contractapi.TransactionContextInterface) (clientID, mspID string, err error) {
	mspID, err = ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return
	}

	cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return
	}

	clientID = cert.Subject.CommonName
	return
}
