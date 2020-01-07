package keys

import (
	"fmt"
	"testing"
)

func TestCreateKeyPair(t *testing.T) {
	kp, err := CreateKeyPair("private.key")
	if err != nil {
		fmt.Println("create key pair error: ", err)
	}
	//skStr, err := kp.GetPrivKeyStr()
	pkStr, err := kp.GetPubKeyStr()
	fmt.Println("pk: ", pkStr)
	fmt.Println("pk: ", len(pkStr))
	//kp.PK.Verify()
}

func TestVerifyIdentifierWithPk(t *testing.T) {
	pkStr := "CAASXjBcMA0GCSqGSIb3DQEBAQUAA0sAMEgCQQDvdIT4CPDRnm8LoWaQkV1QSivBrJuG5J9nI4dG4xhOnzzKf79zzMAYeMylXYpozROwHOmXnmyDaF"
	id := "86.5000.470/0GCSqGSIb3_repo.Mac"
	fmt.Println(VerifyIdentifierWithPk(id, pkStr))
}
