package main

import (
	// "crypto/ecdsa"

	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	rpcClient, err := rpc.Dial("http://172.17.0.5:8575")
	if err != nil {
		// 에러 처리
		logger.LogFatal("Dial : " + err.Error())
		return
	}
	client := ethclient.NewClient(rpcClient)


}

func TempGenerateAdminKeyToPrivateKey() {
	// 보낼 계정의 키 파일 경로를 지정합니다.
	keyFilePath := "./adminKeyFile"
	privateKeyFilePath := "./adminPrivateKey"
	password := "windowshyun"

	// 키 파일과 암호를 사용하여 계정을 언락합니다.
	keyJSON, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		logger.LogFatal("Failed to read key file : " + err.Error())
	}

	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		logger.LogFatal("Failed to decrypt key : " + err.Error())
	}

	// 키를 16진수 문자열로 변환
	hexKey := hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))

	// 16진수 문자열로 저장
	err = ioutil.WriteFile(privateKeyFilePath, []byte(hexKey), 0644)
	if err != nil {
		logger.LogFatal("Failed to write hex key file : " + err.Error())
	}
}

func SendTransaction(client *ethclient.Client, to common.Address, amount float64) (common.Hash, error) {
	// 보낼 계정의 키 파일 경로를 지정합니다.
	privateKeyFilePath := "./adminPrivateKey"

	// 16진수 문자열로 저장된 키 파일을 읽어옴
	hexKeyBytes, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		logger.LogFatal("Failed to read hex key file : " + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to ReadFile")
	}

	// 16진수 문자열을 바이트로 디코딩하여 사용할 수 있음
	decodedKey, err := hex.DecodeString(string(hexKeyBytes))
	if err != nil {
		logger.LogFatal("Failed to decode hex key : " + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to DecodeString")
	}

	// 디코딩된 키를 개인 키로 변환
	privateKey, err := crypto.ToECDSA(decodedKey)
	if err != nil {
		logger.LogFatal("Failed to convert decoded key to ECDSA private key : " + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to ToECDSA")
	}

	// 보낼 계정의 주소를 가져옵니다.
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	logger.LogInfo("보낼 계정 주소 :" + fromAddress.Hex())

	// 보낼 계정의 잔액을 확인합니다.
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		logger.LogFatal("Failed to get balance for address : " + fromAddress.Hex() + " err : " + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to BalanceAt")
	}
	logger.LogInfo("보낼 계정 잔액 :" + WeiToBNB(balance).String())
	logger.LogInfo("받을 계정 주소 :" + to.Hex())

	sendAmount := BNBToWei(amount)

	// 트랜잭션 정보를 생성합니다.
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		logger.LogFatal("Failed to retrieve nonce for address : " + fromAddress.Hex() + " err :" + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to PendingNonceAt")
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		logger.LogFatal("Failed to suggest gas price : " + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to SuggestGasPrice")
	}
	gasLimit := uint64(21000) // 기본 가스 한도

	// 전송할 데이터를 설정합니다. (필요하지 않은 경우 생략 가능)
	data := []byte{}

	// 트랜잭션을 생성합니다.
	tx := types.NewTransaction(nonce, to, sendAmount, gasLimit, gasPrice, data)

	// 체인 아이디를 가져옵니다.
	chainID, err := GetNodeChainID(client)
	if err != nil {
		logger.LogFatal("Failed to GetNodeChainID : " + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to GetNodeChainID")
	}

	// 트랜잭션을 서명합니다.
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		logger.LogFatal("Failed to sign transaction : " + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to SignTx")
	}

	// 트랜잭션을 전송합니다.
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		logger.LogFatal("Failed to send transaction : " + err.Error())
		return common.Hash{}, errors.Wrap(err, "Failed to SendTransaction")
	}

	// 트랜잭션 해시 반환
	return signedTx.Hash(), nil
}