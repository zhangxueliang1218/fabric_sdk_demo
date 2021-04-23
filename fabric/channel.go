package main

import (
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// orgAdmin
func CreateChannel(sdk *fabsdk.FabricSDK, resMgmtClient *resmgmt.Client) (fab.TransactionID, error) {
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		fmt.Printf("mspclient.new err:%s \n", err.Error())
		return "", err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		fmt.Printf("mspclient.GetSigningIdentity err:%s \n", err.Error())
		return "", err
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: "", //integration.GetChannelConfigTxPath(channelID + ".tx"),
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(order1))
	if err != nil {
		fmt.Printf("resMgmtClient.SaveChannel err:%s \n", err.Error())
		return "", err
	}
	return txID.TransactionID, nil
}
