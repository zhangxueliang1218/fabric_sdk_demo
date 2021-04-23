

```
#调用链码
// InitLedger
peer chaincode invoke -o $ORDERER_ADDRESS --tls --cafile $ORDERER_CA   -C zxlchannel -n fabcar --peerAddresses $CORE_PEER_ADDRESS --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["InitLedger"]}' 

#查询链码
// QueryAllCars
peer chaincode query -C zxlchannel -n fabcar -c '{"Args":["QueryAllCars"]}'

// CreateCar
peer chaincode invoke -o $ORDERER_ADDRESS --tls --cafile $ORDERER_CA   -C zxlchannel -n fabcar --peerAddresses $CORE_PEER_ADDRESS --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["CreateCar","hd001","changcheng", "h7", "white", "zxl"]}'


//QueryCar
peer chaincode query -C zxlchannel -n fabcar -c '{"Args":["QueryCar","hd001"]}'

//ChangeCarOwner
peer chaincode invoke -o $ORDERER_ADDRESS --tls --cafile $ORDERER_CA   -C zxlchannel -n fabcar --peerAddresses $CORE_PEER_ADDRESS --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE -c '{"Args":["ChangeCarOwner","hd001","zxl01"]}'


```