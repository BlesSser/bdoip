package client

import (
	"DOIP/models"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

const repoID string = "86.5000.470/0GCSqGSIb3_repo.test"

func TestDOClient_Resolve(t *testing.T) {
	dc := CreateDOClient("pku2018)%!*")
	hr, err := dc.Resolve("86.5000.470/0GCSqGSIb3_repo.test")
	if err != nil {
		fmt.Println("resolve error: ", err)
		return
	}
	fmt.Println(hr)
}

func TestDOClient_Hello(t *testing.T) {
	dc := CreateDOClient("pku2018)%!*")
	flag, err := dc.Hello("86.5000.470/0GCSqGSIb3_repo.test")
	if err != nil {
		fmt.Println("hello error: ", err)
		return
	}
	fmt.Println(flag)
}

func TestDOClient_CreateDO(t *testing.T) {
	dc := CreateDOClient("pku2018)%!*")
	eles := []models.Element{
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
			Data:   []byte("test element"),
		},
	}
	do := &models.DO{
		ID:        dc.HandlePrefix + "data.test2",
		Type:      "string",
		Attr:      "",
		Elements:  eles,
		Signature: "",
	}
	resp, err := dc.CreateDO(repoID, do)
	if err != nil {
		logrus.Error("Create DO error: ", err)
		return
	}
	fmt.Println("resp status: ", resp.Status)
	fmt.Println("resp output: ", resp.Output)
}

func TestDOClient_RetrieveDO(t *testing.T) {
	dc := CreateDOClient("pku2018)%!*")
	do, err := dc.RetrieveDO(dc.HandlePrefix + "data.test1")
	if err != nil {
		logrus.Error("retrieve DO error: ", err)
		return
	}
	fmt.Println(do)
}
