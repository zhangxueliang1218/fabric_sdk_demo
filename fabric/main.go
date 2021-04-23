package main

import (
	"fmt"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID = "zxlchannel"
	orgName   = "org1.zxl.com"
	orgAdmin  = "Admin"
	//ordererOrgName = "OrdererOrg"
	ordererOrgName = "example.com"
	peer1          = "peer1.org1.zxl.com"
	order1         = "orderer1.orderer.zxl.com"

	Org1Msp   = "Org1MSP"
	PeerUser1 = "User1"
)
const (
	ccFilePath  = "D:\\mygithub\\fabric_sdk_demo\\testdata\\go\\src\\github.com\\fabcar"
	ccID        = "fabcar"
	cclabel     = "fabcar_0"
	ccSpec_Type = pb.ChaincodeSpec_GOLANG
)
const (
	configName = "./config/config_e2e.yaml"
)

func main() {
	fmt.Println("-------------- test start--------------")

	//1.
	//CreateChannelTest()

	//2. 2021/04/23 15:30:06 success
	//DeployCCViaLifecycleTest()

	//3. 2021/04/23 15:41:28 success
	SendTransactionTest()

	fmt.Println("-------------- test ending --------------")
}

func CreateChannelTest() {
	fmt.Println("-------------- CreateChannelTest start--------------")
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

	fmt.Println("-------------- CreateChannelTest ending --------------")
}

func DeployCCViaLifecycleTest() {
	fmt.Println("-------------- test start--------------")
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
	DeployCCViaLifecycle(orgResMgmt)
	fmt.Println("-------------- test ending --------------")
}

func SendTransactionTest() {
	fmt.Println("-------------- SendTransactionTest start--------------")
	configProvider := config.FromFile(configName)
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		fmt.Printf("fabsdk.New failed:%s \n", err)
		return
	}
	defer sdk.Close()
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(PeerUser1), fabsdk.WithOrg(orgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Printf("channel.New failed:%s \n", err.Error())
		return
	}
	// peer user
	SendTransaction(client)

	fmt.Println("-------------- SendTransactionTest ending --------------")
}
