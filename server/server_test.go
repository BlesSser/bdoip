package server

import (
	. "DOIP/models"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
)

type TServer struct {
	Port string
}

func (ts *TServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World"))
}

func TestServer_ServeHTTP(t *testing.T) {
	ts := &TServer{Port: ":8081"}
	err := http.ListenAndServe(":8080", ts)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}

func TestServer_Register(t *testing.T) {
	ds, err := CreateServer("127.0.0.1", "10001", "86.5000.470/0GCSqGSIb3_repo.Mac", "../keys/private.key")
	if err != nil {
		return
	}
	pk, err := ds.KeyPair.GetPubKeyStr()
	if err != nil {
		logrus.Fatal("get pk string error.")
	}
	hr := &HandleRecord{
		Identifier: ds.ServerIdentifier,
		Author:     pk,
		Repository: "127.0.0.1:10001",
		Timestamp:  "0",
		Type:       "Repository",
	}
	err = ds.KeyPair.SignHandleRecord(hr)
	if err != nil {
		logrus.Fatal("sign handle record error.")
	}
	err = ds.Register(hr)
	if err != nil {
		logrus.Fatal("register error.", err)
	}
}
