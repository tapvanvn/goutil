package goutil

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

//EthersJSSignMessage compatible with ethers.io signMessage
//as pointed at: https://github.com/ethers-io/ethers.js/issues/823#issuecomment-625953096
func EthersJSSignMessage(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	messageBytes := solsha3.SoliditySHA3WithPrefix(message)
	signature, err := crypto.Sign(messageBytes, privateKey)
	if err != nil {
		return nil, err
	}
	v := uint8(int(signature[64])) + 27
	signature[len(signature)-1] = v
	return signature, nil
}
