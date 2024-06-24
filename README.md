## HLF-project

Этот файл описывает шаги для развертывания цепного кода Hyperledger Fabric, создания, обновления, чтения и удаления активов, а также проверки ограничений доступа по MSP ID и пользователям.

### Шаги для развертывания и тестирования

1. **Настройка переменных окружения для Org1:**

```sh
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051
```

2. **Deploy chain-code:**

```sh
./network.sh down
./network.sh up createChannel -c mychannel -s couchdb 
./network.sh deployCC -ccn sacc -ccp ../../HLF-project/ -ccl go -ccv 1
```

3. **Создание актива от имени Org1:**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -n sacc -c '{"Args":["CreateAsset", "{\"id\":\"asset113\", \"appraisedValue\":1300, \"color\":\"yellow\", \"size\":5, \"type\":\"AssetType\", \"owner\":\"Tom@Org1MSP\"}"]}'
```

4. **Чтение актива от имени Org1:**

```sh
peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -n sacc -c '{"Args":["ReadAsset","asset113"]}'
```

5. **Настройка переменных окружения для Org2:**

```sh
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```

6. **Попытка чтения актива от имени Org2 (должно завершиться неудачно):**

```sh
peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -n sacc -c '{"Args":["ReadAsset","asset113"]}' || echo "Access denied as expected."
```

7. **Попытка обновления актива от имени Org2 (должно завершиться неудачно):**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -n sacc -c '{"Args":["UpdateAsset", "{\"id\":\"asset113\", \"appraisedValue\":1500, \"color\":\"blue\", \"size\":10, \"type\":\"AssetType\", \"owner\":\"Tom@Org1MSP\"}"]}' || echo "Access denied as expected."
```

8. **Попытка удаления актива от имени Org2 (должно завершиться неудачно):**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -n sacc -c '{"Args":["DeleteAsset","asset113"]}' || echo "Access denied as expected."
```

9. **Настройка переменных окружения для Org1:**

```sh
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

10. **Попытка обновления актива от имени Org1 (должно завершиться успешно):**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -n sacc -c '{"Args":["UpdateAsset", "{\"id\":\"asset113\", \"appraisedValue\":1500, \"color\":\"blue\", \"size\":10, \"type\":\"AssetType\", \"owner\":\"Tom@Org1MSP\"}"]}'
```

11. **Попытка удаления актива от имени Org1 (должно завершиться успешно):**

```sh
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -n sacc -c '{"Args":["DeleteAsset","asset113"]}'
```
