/**
 * @Author zxl
 * @Date 2021/4/22 15:49
 * @Desc
 **/
package main

import (
	"fmt"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"github.com/zhangxueliang1218/fabric_sdk_demo/pkg/util"

	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
)

func DeployCCViaLifecycle(orgResMgmt *resmgmt.Client) {
	//Package cc
	ccPkg, err := packageCC(ccFilePath, cclabel, ccSpec_Type)
	if err != nil {
		fmt.Printf("packageCC failed:%s \n", err)
		return
	}

	// Install cc
	packageID, installCCResp, err := installCC(cclabel, ccPkg, orgResMgmt)
	if err != nil {
		fmt.Printf("installCC failed:%s \n", err)
		return
	} else {
		fmt.Printf("Target:%s, Status:%s, PackageID:%s \n", installCCResp.Target, installCCResp.Status, installCCResp.PackageID)
	}

	// Get installed cc package
	ccPkgResp, err := getInstalledCCPackage(packageID, orgResMgmt)
	areEqual := util.ObjectsAreEqual(ccPkg, ccPkgResp)
	if !areEqual {
		fmt.Println("getInstalledCCPackage compare result: false")
		//return
	}

	//Query installed cc
	isInstalled, err := queryInstalled(cclabel, packageID, orgResMgmt)
	if err != nil {
		fmt.Printf("queryInstalled failed:%s \n", err)
		return
	} else {
		fmt.Printf("queryInstalled result: %s \n", isInstalled)
	}

	// Approve cc
	txnID, err := approveCC(ccID, packageID, orgResMgmt)
	if err != nil {
		fmt.Printf("approveCC failed:%s \n", err)
		return
	} else {
		fmt.Printf("approveCC txnID: %s \n", txnID)
	}

	// Query approve cc
	approvedCCResp, err := queryApprovedCC(ccID, orgResMgmt)
	if err != nil {
		fmt.Printf("queryApprovedCC failed:%s \n", err)
		return
	} else {
		fmt.Printf("queryApprovedCC result: %#v \n", approvedCCResp)
	}

	// Check commit readiness
	readiness, err := checkCCCommitReadiness(ccID, orgResMgmt)
	if err != nil {
		fmt.Printf("checkCCCommitReadiness failed:%s \n", err)
		return
	} else {
		fmt.Printf("checkCCCommitReadiness result: %#v \n", readiness)
	}

	//Commit cc
	commitCCtxnID, err := commitCC(ccID, orgResMgmt)
	if err != nil {
		fmt.Printf("commitCC failed:%s \n", err)
		return
	} else {
		fmt.Printf("commitCC result: %#v \n", commitCCtxnID)
	}
	//Query committed cc
	isCommitted, err := queryCommittedCC(ccID, orgResMgmt)
	if err != nil {
		fmt.Printf("queryCommittedCC failed:%s \n", err)
		return
	} else {
		fmt.Printf("queryCommittedCC result: %s \n", isCommitted)
	}
}

func SendTransaction(client *channel.Client) {
	//Init cc
	cc, err := initCC(ccID, client)
	if err != nil {
		fmt.Printf("initCC failed:%s \n", err)
		return
	} else {
		fmt.Printf("initCC result: %#v \n", cc)
	}
	// Invoke cc
	invokeCCResp, err := invokeCC(ccID, client)
	if err != nil {
		fmt.Printf("invokeCC failed:%s \n", err)
		return
	} else {
		fmt.Printf("invokeCC result: %s \n", invokeCCResp.TransactionID)
	}
	// Query cc
	queryCCResp, err := queryCC(ccID, client)
	if err != nil {
		fmt.Printf("queryCC failed:%s \n", err)
		return
	} else {
		fmt.Printf("queryCC result: %s \n", string(queryCCResp.Payload))
	}
}

func packageCC(filePath, label string, ccSpec_Type pb.ChaincodeSpec_Type) ([]byte, error) {
	desc := &lcpackager.Descriptor{
		Path:  filePath,    //integration.GetLcDeployPath_mycc(),//integration.GetLcDeployPath(), //integration.GetLcDeployPath_my(),//
		Type:  ccSpec_Type, // pb.ChaincodeSpec_GOLANG,
		Label: label,       //"mycc_0",//"example_cc_fabtest_e2e_2", //"fabcar_e2e_0",//
	}
	ccPkg, err := lcpackager.NewCCPackage(desc)
	if err != nil {
		fmt.Printf("lcpackager.NewCCPackage failed:%s \n", err.Error())
		return nil, err
	}
	return ccPkg, nil
}

//
// return:
//   packageID, LifecycleInstallCCResponse ,error
func installCC(label string, ccPkg []byte, orgResMgmt *resmgmt.Client) (string, resmgmt.LifecycleInstallCCResponse, error) {
	installCCReq := resmgmt.LifecycleInstallCCRequest{
		Label:   label,
		Package: ccPkg,
	}

	packageID := lcpackager.ComputePackageID(installCCReq.Label, installCCReq.Package)

	resp, err := orgResMgmt.LifecycleInstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Printf("orgResMgmt.LifecycleInstallCC failed:%s \n", err.Error())
		return "", resmgmt.LifecycleInstallCCResponse{}, err
	}
	return packageID, resp[0], nil
}

func getInstalledCCPackage(packageID string, orgResMgmt *resmgmt.Client) ([]byte, error) {
	ccPkg, err := orgResMgmt.LifecycleGetInstalledCCPackage(packageID, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Printf("orgResMgmt.LifecycleGetInstalledCCPackage failed:%s \n", err.Error())
		return nil, err
	}
	return ccPkg, nil
}

func queryInstalled(label string, packageID string, orgResMgmt *resmgmt.Client) (bool, error) {
	resp, err := orgResMgmt.LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(peer1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Printf("orgResMgmt.LifecycleQueryInstalledCC failed:%s \n", err.Error())
		return false, err
	}
	for _, cc := range resp {
		if cc.PackageID == packageID && cc.Label == label {
			return true, nil
		}
	}
	return false, nil
}

func approveCC(ccID, packageID string, orgResMgmt *resmgmt.Client) (fab.TransactionID, error) {
	ccPolicy := policydsl.SignedByAnyMember([]string{Org1Msp})
	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              ccID,
		Version:           "0",
		PackageID:         packageID,
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}

	txnID, err := orgResMgmt.LifecycleApproveCC(channelID, approveCCReq, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithOrdererEndpoint(order1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Printf("orgResMgmt.LifecycleApproveCC failed:%s \n", err.Error())
		return "", err
	}
	return txnID, nil
}

func queryApprovedCC(ccID string, orgResMgmt *resmgmt.Client) (resmgmt.LifecycleApprovedChaincodeDefinition, error) {
	queryApprovedCCReq := resmgmt.LifecycleQueryApprovedCCRequest{
		Name:     ccID,
		Sequence: 1,
	}
	resp, err := orgResMgmt.LifecycleQueryApprovedCC(channelID, queryApprovedCCReq, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Printf("orgResMgmt.LifecycleQueryApprovedCC failed:%s \n", err.Error())
		return resmgmt.LifecycleApprovedChaincodeDefinition{}, err
	}
	return resp, nil
}

func checkCCCommitReadiness(ccID string, orgResMgmt *resmgmt.Client) (resmgmt.LifecycleCheckCCCommitReadinessResponse, error) {
	ccPolicy := policydsl.SignedByAnyMember([]string{Org1Msp})
	req := resmgmt.LifecycleCheckCCCommitReadinessRequest{
		Name:              ccID,
		Version:           "0",
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		Sequence:          1,
		InitRequired:      true,
	}
	resp, err := orgResMgmt.LifecycleCheckCCCommitReadiness(channelID, req, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Printf("orgResMgmt.LifecycleCheckCCCommitReadiness failed:%s \n", err.Error())
		return resmgmt.LifecycleCheckCCCommitReadinessResponse{}, err
	}
	return resp, nil
}

func commitCC(ccID string, orgResMgmt *resmgmt.Client) (fab.TransactionID, error) {
	ccPolicy := policydsl.SignedByAnyMember([]string{Org1Msp})
	req := resmgmt.LifecycleCommitCCRequest{
		Name:              ccID,
		Version:           "0",
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	txnID, err := orgResMgmt.LifecycleCommitCC(channelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(peer1), resmgmt.WithOrdererEndpoint(order1))
	if err != nil {
		fmt.Printf("orgResMgmt.LifecycleCheckCCCommitReadiness failed:%s \n", err.Error())
		return "", err
	}
	return txnID, nil
}

func queryCommittedCC(ccID string, orgResMgmt *resmgmt.Client) (bool, error) {
	req := resmgmt.LifecycleQueryCommittedCCRequest{
		Name: ccID,
	}
	resp, err := orgResMgmt.LifecycleQueryCommittedCC(channelID, req, resmgmt.WithTargetEndpoints(peer1), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Printf("orgResMgmt.LifecycleCheckCCCommitReadiness failed:%s \n", err.Error())
		return false, err
	}
	for _, cc := range resp {
		if cc.Name == ccID {
			return true, nil
		}
	}
	return false, nil
}

func initCC(ccID string, client *channel.Client) (channel.Response, error) {
	args := [][]byte{}
	fcn := "InitLedger"
	isInit := true
	resp, err := client.Execute(channel.Request{ChaincodeID: ccID, Fcn: fcn, Args: args, IsInit: isInit}, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		fmt.Printf("client.Execute failed:%s \n", err.Error())
		return channel.Response{}, err
	}
	return resp, nil
}

func invokeCC(ccID string, client *channel.Client) (channel.Response, error) {
	var args [][]byte
	args = append(args, []byte("hd001"))
	args = append(args, []byte("changcheng"))
	args = append(args, []byte("h7"))
	args = append(args, []byte("white"))
	args = append(args, []byte("zxl"))

	invkerRequest := channel.Request{
		ChaincodeID: ccID,
		Fcn:         "CreateCar",
		Args:        args,
	}
	resp, err := client.Execute(invkerRequest) //handler
	if err != nil {
		fmt.Printf("channelClient.Execute failed:%s \n", err)
		return channel.Response{}, err
	}
	return resp, nil
}

func queryCC(ccID string, client *channel.Client) (channel.Response, error) {
	var args [][]byte
	args = append(args, []byte("hd001"))
	//
	request := channel.Request{
		ChaincodeID: ccID,
		Fcn:         "QueryCar",
		Args:        args,
	}
	resp, err := client.Query(request)
	if err != nil {
		fmt.Printf("channelClient.Query failed:%s \n", err)
		return channel.Response{}, err
	}
	return resp, nil
}
