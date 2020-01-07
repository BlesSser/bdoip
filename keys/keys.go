package keys

import (
	. "DOIP/models"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type KeyPairs struct {
	PK crypto.PubKey
	SK crypto.PrivKey
}

func CreateKeyPair(skFile string) (*KeyPairs, error) {

	sk, pk, err := loadKeys(skFile)
	if err != nil {
		logrus.Error("load Keys error", err)
		return nil, err
	}

	return &KeyPairs{
		PK: pk,
		SK: sk,
	}, nil
}

//loadKeys load sk and generate pk from given file
func loadKeys(skFile string) (sk crypto.PrivKey, pk crypto.PubKey, err error) {
	skb, err := ioutil.ReadFile(skFile)
	if err != nil {
		//generate new key pair
		if os.IsNotExist(err) {
			sk, pk, err1 := crypto.GenerateKeyPair(crypto.RSA, 512)
			if err1 != nil {
				logrus.Error("generate key pair error: ", err1)
				return nil, nil, err1
			}
			skStr, err1 := keyToString(sk)
			if err1 != nil {
				logrus.Error("private key to string error: ", err1)
				return nil, nil, err1
			}
			err1 = ioutil.WriteFile(skFile, []byte(skStr), os.ModePerm)
			if err1 != nil {
				logrus.Error("write Key to skFile error: ", err1)
				return nil, nil, err1
			}
			return sk, pk, nil
		} else {
			fmt.Println("read sk error : ", err)
			return
		}

	}

	skb2, err := crypto.ConfigDecodeKey(string(skb))
	if err != nil {
		fmt.Println("decode sk error : ", err)
		return
	}

	sk, err = crypto.UnmarshalPrivateKey(skb2)
	if err != nil {
		fmt.Println("unmarshal sk error : ", err)
		return
	}
	pk = sk.GetPublic()

	return
}

//GetPubKeyStr format public key to string
func (kp *KeyPairs) GetPubKeyStr() (keyStr string, err error) {
	keyStr, err = keyToString(kp.PK)
	if err != nil {
		fmt.Println("public key to Bytes error : ", err)
		return
	}
	return keyStr, nil
}

//GetPrivKeyStr format private key to string
func (kp *KeyPairs) GetPrivKeyStr() (keyStr string, err error) {
	keyStr, err = keyToString(kp.SK)
	if err != nil {
		fmt.Println("private key to Bytes error : ", err)
		return
	}
	return keyStr, nil
}

func keyToString(key crypto.Key) (keyStr string, err error) {
	keyB, err := key.Bytes()
	if err != nil {
		fmt.Println("key to Bytes error : ", err)
		return
	}
	keyStr = crypto.ConfigEncodeKey(keyB)
	return keyStr, nil
}

func (kp *KeyPairs) SignHandleRecord(hr *HandleRecord) (err error) {
	body := Concat(
		[]byte(hr.Author),
		[]byte(hr.Identifier),
		[]byte(hr.Repository),
		[]byte(hr.Timestamp),
		[]byte(hr.Type),
	)
	sign, err := kp.SK.Sign(body)
	if err != nil {
		logrus.Println("sign HandleRecord body error")
	}
	signStr := crypto.ConfigEncodeKey(sign)
	hr.Signature = signStr
	return
}

func (kp *KeyPairs) SignResponse(resp *Response) (err error) {
	body := Concat(
		[]byte(resp.Attr),
		[]byte(resp.Output),
		[]byte(resp.PKString),
		[]byte(resp.RequestID),
		[]byte(resp.Status),
	)
	sign, err := kp.SK.Sign(body)
	if err != nil {
		logrus.Println("sign HandleRecord body error")
	}
	signStr := crypto.ConfigEncodeKey(sign)
	resp.Signature = signStr
	return
}
