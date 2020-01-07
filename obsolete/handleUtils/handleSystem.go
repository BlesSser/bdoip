package handleUtils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
)

const HandleParseUrl = "http://39.104.122.7:9000/handleSystem/service/HandleParseService?wsdl"
const HandleRegisterUrl = "https://39.104.122.7:9443/handleSystem/service/Register/registerService?wsdl"
const HandlePrefix = "86.5000.470/"

//ResolveHandle resolve handle to an RecordList
func ResolveHandle(handleId ...string) (result string, err error) {
	handleStr := ""
	for _, str := range handleId {
		handleStr += " <handleId>" + str + "</handleId>"
	}
	fmt.Println("handleID : \n", handleStr)
	soapRequestData := "<?xml version=\"1.0\" encoding=\"utf-8\"?>" +
		"<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"<soap:Body>" + "<ns1:handleParse xmlns:ns1=\"http://service.handleparse.cdi.cn/\">" +
		handleStr +
		"  <isAuthentication>0</isAuthentication>" +
		"  <relationFlag>1</relationFlag>" +
		"  <deepRelationFlag>1</deepRelationFlag>" +
		"  <isReverse>0</isReverse>" +
		"  <userGroupIds></userGroupIds>" +
		"  <onlyHandleIds>0</onlyHandleIds>" +
		"  <appHandleId></appHandleId>" +
		"  <appHandleIdEncryptBase64></appHandleIdEncryptBase64>" +
		"  <appHandleIdSignBase64></appHandleIdSignBase64>" +
		"  <password></password>" +
		"</ns1:handleParse>" + " </soap:Body>" + "</soap:Envelope>"
	fmt.Println(soapRequestData)
	res, err := http.Post(HandleParseUrl, "application/soap+xml;charset=utf-8", bytes.NewBuffer([]byte(soapRequestData)))
	if err != nil {
		fmt.Println("Fatal error while get from handle", err.Error())
		return
	}

	bd, _ := ioutil.ReadAll(res.Body)
	result = html.UnescapeString(string(bd))
	return
}

//解析返回的XML字符串，提取有效的Record
func ParseResponseToRecord(res string) Records {

	fmt.Println("\nResponse in xml : \n", res)

	root := &MyRespEnvelope{}
	err := xml.Unmarshal([]byte(res), &root)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return root.Body.HandleParseResponse.Return.Root.ListRecords
}

//向Handle系统中注册resourceObject
func RegisterResource() (result string) {

	handleStr :=
		"<root xmlns=\"http://www.cdi.cn/handle\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">" +
			"<ListRecords metedataStandard=\"86.5000.470/metadatastandard.resourceObject\">" +
			"<record dataLevel=\"1\" regType=\"1\">" +
			"<header>" +
			"<identifier>86.5000.470/resourceObjectTest</identifier>" +
			"<metadataOption>1</metadataOption>" +
			"</header>" +
			"<metadata>" +
			"<repository><value>" +
			"<![CDATA[127.0.0.1:8080]]>" +
			"</value></repository>" +
			"<type><value>" +
			"<![CDATA[resourceObject]]>" +
			"</value></type>" +
			"<description><value>" +
			"<![CDATA[测试注册接口]]>" +
			"</value></description>" +
			"<timestamp><value>" +
			"<![CDATA[0]]>" +
			"</value></timestamp>" +
			"<CHS_DUPLICATECHECK><value>" +
			"<![CDATA[0]]>" +
			"</value></CHS_DUPLICATECHECK>" +
			"</metadata>" +
			"</record>" +
			"</ListRecords>" +
			"</root>"
	handleStr = ""
	//fmt.Println("handleStr : \n",handleStr)
	soapRequestData := "<?xml version=\"1.0\" encoding=\"utf-8\"?>" +
		"<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"<soap:Body>" + "<ns1:handleRegister xmlns:ns1=\"http://service.handleRegistered.cdi.cn/\">" +
		"<xml>" + handleStr + "</xml>" +
		"<adminHandleId>0.NA/86.5000.470</adminHandleId>" +
		"<adminIndex>1</adminIndex>" +
		"<applicationHandleId></applicationHandleId>" +
		"<appHandleIdEncryptBase64></appHandleIdEncryptBase64>" +
		"<appHandleIdSignBase64></appHandleIdSignBase64>" +
		"<password></password>" +
		"<operationTime></operationTime>" +
		"</ns1:handleRegister>" + " </soap:Body>" + "</soap:Envelope>"

	fmt.Println(soapRequestData)
	//x509.Certificate.
	pool := x509.NewCertPool()
	//caCertPath := "etcdcerts/ca.crt"
	caCertPath := "../keys/86.5000.470private.pfx"
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)
	//pool.AddCert(caCrt)

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	res, err := client.Post(HandleRegisterUrl, "application/soap+xml;charset=utf-8", bytes.NewBuffer([]byte(soapRequestData)))
	if err != nil {
		fmt.Println("Fatal error while get from handle", err.Error())
		return
	}

	bd, _ := ioutil.ReadAll(res.Body)
	result = html.UnescapeString(string(bd))
	return

}

func RegisterWithLib() (err error) {
	//soap, err := gosoap.SoapClient(HandleParseUrl)
	//if err != nil {
	//	log.Fatalf("SoapClient error: %s", err)
	//}
	//
	//params := gosoap.Params{
	//	"handleId": "86.5000.470/haier",
	//	"isAuthentication":"0",
	//	"relationFlag":"1",
	//	"deepRelationFlag":"1",
	//	"isReverse":"0",
	//	"userGroupIds":"",
	//	"onlyHandleIds":"0",
	//	"appHandleId":"",
	//	"appHandleIdEncryptBase64":"",
	//	"appHandleIdSignBase64":"",
	//	"password":"",
	//}
	//res, err := soap.Call("handleParse", params)
	//if err != nil {
	//	log.Fatalf("Call error: %s", err)
	//}
	//fmt.Println(soap.HeaderParams)
	////root := &MyRespEnvelope{}
	//
	////err = res.Unmarshal(&root)
	////if err != nil{
	////	return err
	////}
	//fmt.Println(string(res.Body))
	return
}

func GetData(url string) (result string) {
	//跳过https证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Transport: tr,
	}
	response, err := c.Get(url)
	if err != nil {
		fmt.Println("Error while get ApiData: ", err)
		return
	} else {
		bd, _ := ioutil.ReadAll(response.Body)
		result = html.UnescapeString(string(bd))
		return
	}

}
