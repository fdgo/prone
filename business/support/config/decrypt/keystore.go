package decrypt

import (
	"business/support/libraries/loggers"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func KAuth(addressStr string, signatureStr string, body []byte) error {
	if len(body) <= 0 {
		return errors.New("no sign")
	}
	if ok, err := verifySignature(addressStr, signatureStr, body); !ok {
		if nil != err {
			return err
		} else {
			return errors.New("verify signature failed!")
		}
	}
	return nil
}

func KSign(message []byte, privateKey string) (string, error) {
	if nil == message ||
		len(message) <= 0 ||
		len(privateKey) <= 0 {
		return "", errors.New("invalid params")
	}
	hexKey, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	if nil == hexKey || len(hexKey) <= 0 {
		return "", errors.New("invalid params")
	}
	pKey, err := crypto.ToECDSA(hexKey)
	if err != nil {
		return "", err
	}
	if nil == pKey {
		return "", errors.New("invalid params")
	}
	signature, err := crypto.Sign(signHash(message), pKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

func verifySignature(
	addressStr string,
	signatureHexStr string,
	message []byte) (bool, error) {
	// Determine whether it is a keyfile, public key or address.
	if !common.IsHexAddress(addressStr) {
		return false, errors.New("Invalid address")
	}
	address := common.HexToAddress(addressStr)

	signature, err := hex.DecodeString(signatureHexStr)
	if err != nil {
		return false, err
	}

	if len(message) == 0 {
		return false, errors.New("A message must be provided")
	}
	// Read message if file.
	// if _, err := os.Stat(string(message)); err == nil {
	// 	message, err = ioutil.ReadFile(string(message))
	// 	if err != nil {
	// 		loggers.Warn.Printf("Failed to read the message file: %v", err)
	// 		return false, err
	// 	}
	// }

	recoveredPubkey, err := crypto.SigToPub(signHash(message), signature)
	if err != nil || recoveredPubkey == nil {
		loggers.Warn.Printf("Signature verification failed: %v,recoveredPubkey:%v", err, recoveredPubkey)
		return false, err
	}
	// recoveredPubkeyBytes := crypto.FromECDSAPub(recoveredPubkey)
	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	success := address == recoveredAddress
	return success, nil
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
