package handleUtils

import (
	"fmt"
	"testing"
)

func TestResolveHandle(t *testing.T) {
	handle := "86.5000.470/haier"
	res, err := ResolveHandle(handle)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	records := ParseResponseToRecord(res)
	fmt.Println("RESULT: ", records.Record[0].MetaData.ApiAddress.Value)
}

func TestRegisterResource(t *testing.T) {
	res := RegisterResource()
	fmt.Println(res)
}

func TestRegisterWithLib(t *testing.T) {
	err := RegisterWithLib()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
