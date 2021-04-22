package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	configName = "./myfabric/config_e2e_testenv.yaml"
)

func main() {
	fmt.Println("test")
	configProvider := config.FromFile(configName)
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		fmt.Printf("fabsdk.New failed:%s \n", err)
		return
	}
	defer sdk.Close()
	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		fmt.Printf("Failed to create new resource management client: %s \n", err)
		return
	}
	// orgAdmin
	CreateChannel(sdk, orgResMgmt)

	// orgAdmin
	DeployCCViaLifecycleTest(orgResMgmt)

	// peer user
	SendTransactionTest(sdk)
}
