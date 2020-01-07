package client

import (
	"DOIP/keys"
	"DOIP/models"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DOClient struct {
	Auth         string
	ResolveAddr  string
	RepoMap      map[string]string
	HandlePrefix string
}

const ResolvutionURL string = "http://47.106.38.23:10001/resolve?identifier="
const HandlePrefix string = "86.5000.470/0GCSqGSIb3_"

func CreateDOClient(auth string) *DOClient {
	return &DOClient{
		Auth:         auth,
		ResolveAddr:  ResolvutionURL,
		RepoMap:      make(map[string]string),
		HandlePrefix: HandlePrefix,
	}
}

//Resolve 将标识解析为Handle Record，并且验证其合法性
func (dc *DOClient) Resolve(identifier string) (hr *models.HandleRecord, err error) {
	resp, err := http.Get(dc.ResolveAddr + identifier)
	if err != nil {
		logrus.Error("connect to handle server error: ", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	hr = &models.HandleRecord{}
	err = json.Unmarshal(body, hr)
	if err != nil {
		logrus.Error("parse response to handle record error: ", err)
		return
	}
	flag, err := keys.VerifyHandleRecord(hr)
	if err != nil {
		logrus.Error("verify handle record error: ", err)
		return
	}
	if flag {
		return
	}
	err = errors.New("invalid handle record")
	return
}

//Hello 解析一个RepositoryID，然后向其发送DOIP请求，查看其是否支持DOIP协议
func (dc *DOClient) Hello(identifier string) (out string, err error) {
	repoAddr, err := dc.getRepoAddr(identifier)
	if err != nil {
		logrus.Error("get repository address error: ", identifier)
		return
	}
	req := models.Request{
		TargetID:       identifier,
		RequestID:      "1",
		ClientID:       dc.Auth,
		OperationID:    models.OPERATION_HELLO,
		Attr:           "",
		Authentication: dc.Auth,
		Input:          "",
	}
	resp, err := http.PostForm(repoAddr, url.Values{
		"TargetID":       {req.TargetID},
		"RequestID":      {req.RequestID},
		"ClientID":       {req.ClientID},
		"OperationID":    {string(req.OperationID)},
		"Attr":           {req.Attr},
		"Authentication": {req.Authentication},
		"Input":          {req.Input},
	})
	if err != nil {
		logrus.Errorf("hello to repository error:%v, repo addr:%v \n", err, repoAddr)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("read response body error:%v, repo addr:%v \n", err, repoAddr)
		return
	}
	doipResp := &models.Response{}
	err = json.Unmarshal(body, doipResp)
	if err != nil {
		logrus.Errorf("parse json to doip response error:%v, response:%v \n ", err, string(body))
		return
	}
	if doipResp.Status != models.STATUS_SUCESS {
		return
	}

	return doipResp.Output, nil
}

//RetrieveDO 根据标识，从DORepository 获取DO
func (dc *DOClient) RetrieveDO(identifier string) (do *models.DO, err error) {
	hr, err := dc.Resolve(identifier)
	if err != nil {
		logrus.Errorf("get handle record error:%v, response:%v \n ", err, identifier)
		return
	}
	logrus.Info("get DO from Repository: ", hr.Repository)
	repoAddr, err := dc.getRepoAddr(hr.Repository)
	logrus.Info("repoAddr: ", repoAddr)
	if err != nil {
		logrus.Error("get repository address error: ", identifier)
		return
	}
	req := &models.Request{
		TargetID:       identifier,
		RequestID:      "1",
		ClientID:       dc.Auth,
		OperationID:    models.OPERATION_RETRIEVE,
		Attr:           "includeElement",
		Authentication: dc.Auth,
		Input:          "",
	}
	doipResp, err := dc.getAndVerifyResponse(req, repoAddr)
	if err != nil {
		logrus.Errorf("get response error:%v \n ", err)
		return
	}
	if doipResp.Status != models.STATUS_SUCESS {
		logrus.Warnf("failed server status : %v, output:%v  \n ", models.STATUS_SUCESS, doipResp.Output)
		return nil, errors.New("failed server status")
	}
	doStr := doipResp.Output
	logrus.Info("output str: ", doStr)
	do = &models.DO{}
	err = json.Unmarshal([]byte(doStr), do)
	if err != nil {
		logrus.Errorf("unmarshal DO error:%v, doStr:%v \n ", err, doipResp.Output)
		return nil, err
	}
	return
}

//CreateDO 在指定Repository上，创建DO
//输入repoID，do，DO签名在Repo上进行
func (dc *DOClient) CreateDO(repoID string, do *models.DO) (resp *models.Response, err error) {
	repoAddr, err := dc.getRepoAddr(repoID)
	if err != nil {
		logrus.Errorf("get repository address error:%v \n ", err)
		return
	}
	doBytes, err := json.Marshal(do)
	if err != nil {
		logrus.Errorf("marshal digital object error:%v \n ", err)
		return
	}
	req := &models.Request{
		TargetID:       do.ID,
		RequestID:      "1",
		ClientID:       dc.Auth,
		OperationID:    models.OPERATION_CREATE,
		Attr:           "",
		Authentication: dc.Auth,
		Input:          string(doBytes),
	}
	resp, err = dc.getAndVerifyResponse(req, repoAddr)
	if err != nil {
		logrus.Errorf("get and verify response error:%v \n ", err)
		return
	}
	if resp.Status != models.STATUS_SUCESS {
		logrus.Warnf("failed server status : %v, output:%v  \n ", models.STATUS_SUCESS, resp.Output)
		return nil, errors.New("failed server status")
	}
	return
}

//DeleteDO 根据标识，从DORepository 获取DO
func (dc *DOClient) DeleteDO() {

}

//ListOperation 根据标识，从DORepository 获取DO
func (dc *DOClient) ListOperation() {

}

func (dc *DOClient) getRepoAddr(identifier string) (repoAddr string, err error) {
	var hr *models.HandleRecord
	if v, ok := dc.RepoMap[identifier]; !ok {
		hr, err = dc.Resolve(identifier)
		if err != nil {
			logrus.Error("resolve repository identifier error: ", err)
			return
		}
		if hr.Type != "Repository" {
			logrus.Error("this is not a repository ID: ", identifier)
			return "", errors.New("not a repository identifier")
		}
		repoAddr = hr.Repository
	} else {
		repoAddr = v
	}
	return
}

//getResponse 向Repository发送DOIP request 获取response并验证其合法性
func (dc *DOClient) getAndVerifyResponse(req *models.Request, repoAddr string) (doipResp *models.Response, err error) {
	resp, err := http.PostForm(repoAddr, url.Values{
		"TargetID":       {req.TargetID},
		"RequestID":      {req.RequestID},
		"ClientID":       {req.ClientID},
		"OperationID":    {string(req.OperationID)},
		"Attr":           {req.Attr},
		"Authentication": {req.Authentication},
		"Input":          {req.Input},
	})
	if err != nil {
		logrus.Errorf("retrieve from repository error:%v, repo addr:%v \n", err, repoAddr)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("read response body error:%v, repo addr:%v \n", err, repoAddr)
		return
	}
	doipResp = &models.Response{}
	err = json.Unmarshal(body, doipResp)
	if err != nil {
		logrus.Errorf("parse json to doip response error:%v, response:%v \n ", err, string(body))
		return
	}
	valid, err := keys.VerifyResponse(doipResp)
	if err != nil {
		logrus.Errorf("verify response error:%v, response:%v \n ", err, string(body))
		return
	}
	if !valid {
		logrus.Errorf("invalid response: %v\n", string(body))
		return nil, errors.New("invalid response error")
	}
	return
}
