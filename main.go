package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://39.104.122.7:9443/handleSystem/service/HandleParseService?wsdl")
	if err != nil {
		fmt.Print("ERROR: ", err)
	}
	fmt.Print("RESP: ", resp)
}
