package models

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestDOtoJSON(t *testing.T) {

	eles := []Element{
		{
			ID:     "eleID1",
			Length: 1,
			Type:   "string",
			Attr:   "1",
			Data:   []byte("bytes"),
		},
		{
			ID:     "eleID2",
			Length: 2,
			Type:   "string",
			Attr:   "2",
			Data:   []byte("123"),
		},
	}
	do := &DO{
		ID:        "id",
		Type:      "",
		Attr:      "",
		Elements:  eles,
		Signature: "",
	}
	fmt.Println(do)
	bytes, err := json.Marshal(do)
	if err != nil {
		fmt.Println("parse json error", err)
	}
	str := string(bytes)
	fmt.Println(str)
	err = json.Unmarshal(bytes, do)
	if err != nil {
		fmt.Println("parse json error", err)
	}
	fmt.Println(do.Elements[0].Data)
}

func TestTrifle(t *testing.T) {
	fmt.Println(time.Now().Second() + time.Now().Nanosecond())
}
