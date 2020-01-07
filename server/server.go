package server

import (
	"DOIP/keys"
	. "DOIP/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Server struct {
	Port                string
	IP                  string
	ServerIdentifier    string
	SKFile              string
	KeyPair             *keys.KeyPairs
	HandleHello         Hello
	HandleCreate        Create
	HandleRetrieve      Retrieve
	HandleUpdate        Update
	HandleDelete        Delete
	HandleSearch        Search
	HandleListOperation ListOperation
}

//创建一个doip服务，传入标识，私钥。
func CreateServer(ip string, port string, ID string, skFile string) (ds *Server, err error) {
	kp, err := keys.CreateKeyPair(skFile)
	if err != nil {
		logrus.Fatal("create key pairs error")
		return
	}
	pkStr, err := kp.GetPubKeyStr()
	logrus.Info("PK string: ", pkStr)
	if !keys.VerifyIdentifierWithPk(ID, pkStr) {
		err = errors.New("identifier and PK does not match")
		return
	}
	ds = &Server{
		Port:                port,
		IP:                  ip,
		ServerIdentifier:    ID,
		SKFile:              skFile,
		KeyPair:             kp,
		HandleHello:         nil,
		HandleCreate:        nil,
		HandleRetrieve:      nil,
		HandleUpdate:        nil,
		HandleDelete:        nil,
		HandleSearch:        nil,
		HandleListOperation: nil,
	}
	return
}

//启动DOIP服务，自动向服务器端注册本地服务
func (ds *Server) Start() (err error) {

	pkStr, err := ds.KeyPair.GetPubKeyStr()

	ts := time.Now().UnixNano()
	hr := &HandleRecord{
		Identifier: ds.ServerIdentifier,
		Author:     pkStr,
		Repository: "http://" + ds.IP + ds.Port,
		Timestamp:  strconv.FormatInt(ts, 10),
		Type:       "Repository",
		Signature:  "",
	}
	err = ds.KeyPair.SignHandleRecord(hr)
	if err != nil {
		logrus.Fatal("sign before register error: ", err)
		return
	}
	err = ds.Register(hr)
	if err != nil {
		logrus.Fatal("register repository error: ", err)
		return
	}
	return http.ListenAndServe(ds.Port, ds)
}

//注册一条Handle Record
func (ds *Server) Register(hr *HandleRecord) (err error) {
	if !keys.VerifyIdentifierWithPk(hr.Identifier, hr.Author) {
		logrus.Error("identifier does not match author")
		return errors.New("identifier does not match author")
	}
	flag, err := keys.VerifyHandleRecord(hr)
	if err != nil {
		logrus.Error("Verify Handle Record Error")
		return
	}
	if !flag {
		logrus.Error("signature does not match pk")
		return errors.New("signature does not match pk")
	}
	resp, err := http.PostForm(RegisterUrl, url.Values{
		"identifier": {hr.Identifier},
		"repository": {hr.Repository},
		"type":       {hr.Type},
		"author":     {hr.Author},
		"timestamp":  {hr.Timestamp},
		"signature":  {hr.Signature},
	})
	fmt.Println("signature: ", hr.Signature)
	if err != nil {
		logrus.Error("register post request error: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("got response error: ", err)
		return
	}
	logrus.Infof("successfully register DO:%v, response:%v ", hr.Identifier, string(body))
	return
}

type Hello func() (output string, status STATUS)
type Create func(do *DO, auth string) (output string, status STATUS, err error)
type Retrieve func(handle string, includeElement bool) (do *DO, status STATUS, err error)
type Update func(do *DO, auth string) (output string, status STATUS, err error)
type Delete func(handle string, auth string) (output string, status STATUS, err error)
type Search func(query string) (output string, status STATUS, err error)
type ListOperation func(handle string) (output string, status STATUS, err error)

func (ds *Server) SetHello(f Hello) {
	ds.HandleHello = f
}

func (ds *Server) SetCreate(f Create) {
	ds.HandleCreate = f
}

func (ds *Server) SetRetrieve(f Retrieve) {
	ds.HandleRetrieve = f
}

func (ds *Server) SetUpdate(f Update) {
	ds.HandleUpdate = f
}
func (ds *Server) SetDelete(f Delete) {
	ds.HandleDelete = f
}

func (ds *Server) SetSearch(f Search) {
	ds.HandleSearch = f
}

func (ds *Server) SetListOperation(f ListOperation) {
	ds.HandleListOperation = f
}

//ServeHTTP，在HTTP协议的基础上实现DOIP协议
//部分实现与DOIP不同，request类型较少，且一个连接一个request，DOIP中一个连接可以包含多个request
func (ds *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logrus.Error("parse DOI error: ", err)
		_, _ = w.Write([]byte("parse DOI error, check parameters"))
		return
	}
	request := &Request{}
	if len(r.Form["TargetID"]) == 0 {
		_, _ = w.Write([]byte("Target ID not found, check parameters"))
		return
	}
	request.TargetID = r.Form["TargetID"][0]
	if len(r.Form["OperationID"]) == 0 {
		_, _ = w.Write([]byte("Operation ID not found, check parameters"))
		return
	}
	request.OperationID = OPERATION(r.Form["OperationID"][0])
	if len(r.Form["Attr"]) != 0 {
		request.Attr = r.Form["Attr"][0]
	}
	if len(r.Form["ClientID"]) != 0 {
		request.ClientID = r.Form["ClientID"][0]
	}
	if len(r.Form["RequestID"]) != 0 {
		request.RequestID = r.Form["RequestID"][0]
	}
	if len(r.Form["Authentication"]) != 0 {
		request.Authentication = r.Form["Authentication"][0]
	}
	if len(r.Form["Input"]) != 0 {
		request.Input = r.Form["Input"][0]
	}
	pkStr, err := ds.KeyPair.GetPubKeyStr()
	response := &Response{
		RequestID: request.RequestID,
		PKString:  pkStr,
		Status:    STATUS_SUCESS,
	}
	switch request.OperationID {
	case OPERATION_HELLO:
		if ds.HandleHello != nil {
			response.Output, response.Status = ds.HandleHello()
		} else {
			response.Status = STATUS_REQUEST_INVALID
			response.Output = "DOIPServer do not support " + string(request.OperationID)
		}
	case OPERATION_CREATE:
		if ds.HandleCreate != nil {
			do := &DO{}
			err := json.Unmarshal([]byte(request.Input), do)
			if err != nil {
				response.Status = STATUS_REQUEST_INVALID
				response.Output = "Parse JSON to DO error, Please check DO string"
			}
			response.Output, response.Status, err = ds.HandleCreate(do, request.Authentication)
			if err != nil {
				logrus.Error(err)
			}
		} else {
			response.Status = STATUS_REQUEST_INVALID
			response.Output = "DOIPServer do not support " + string(request.OperationID)
		}
	case OPERATION_RETRIEVE:
		includeElement := false
		if request.Attr == "includeElement" {
			includeElement = true
		}
		if ds.HandleRetrieve != nil {
			output, status, err := ds.HandleRetrieve(request.TargetID, includeElement)
			response.Status = status
			if err != nil {
				logrus.Error(err)
			} else {
				if output != nil {
					outputByte, _ := json.Marshal(output)
					response.Output = string(outputByte)
				}
			}
		} else {
			response.Status = STATUS_REQUEST_INVALID
			response.Output = "DOIPServer do not support " + string(request.OperationID)
		}
	case OPERATION_UPDATE:
		if ds.HandleUpdate != nil {
			do := &DO{}
			err := json.Unmarshal([]byte(request.Input), do)
			if err != nil {
				response.Status = STATUS_REQUEST_INVALID
				response.Output = "Parse JSON to DO error, Please check DO string"
				logrus.Error(err)
			}
			output, status, err := ds.HandleUpdate(do, request.Authentication)
			response.Output = output
			response.Status = status
			if err != nil {
				logrus.Error(err)
			}
		} else {
			response.Status = STATUS_REQUEST_INVALID
			response.Output = "DOIPServer do not support " + string(request.OperationID)
		}
	case OPERATION_DELETE:
		if ds.HandleDelete != nil {
			response.Output, response.Status, err = ds.HandleDelete(request.TargetID, request.ClientID)
			if err != nil {
				logrus.Error(err)
			}
		} else {
			response.Status = STATUS_REQUEST_INVALID
			response.Output = "DOIPServer do not support " + string(request.OperationID)
		}
	case OPERATION_SEARCH:
		if ds.HandleSearch != nil {
			response.Output, response.Status, err = ds.HandleSearch(request.Input)
			if err != nil {
				logrus.Error(err)
			}
		} else {
			response.Status = STATUS_REQUEST_INVALID
			response.Output = "DOIPServer do not support " + string(request.OperationID)
		}
	case OPERATION_LISTOPERATIONS:
		if ds.HandleListOperation != nil {
			response.Output, response.Status, err = ds.HandleListOperation(request.TargetID)
			if err != nil {
				logrus.Error(err)
			}
		} else {
			response.Status = STATUS_REQUEST_INVALID
			response.Output = "DOIPServer do not support " + string(request.OperationID)
		}
	default:
		response.Status = STATUS_REQUEST_INVALID
		response.Output = "Invalid operation ID " + string(request.OperationID)
	}
	err = ds.KeyPair.SignResponse(response)
	if err != nil {
		logrus.Error("sign response error")
		response.Status = STATUS_MULTI_ERRORS
	}
	httpOut, err := json.Marshal(response)
	if err != nil {

	}
	_, _ = w.Write(httpOut)
}
