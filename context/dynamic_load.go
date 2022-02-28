package context

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/free5gc/smf/factory"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/free5gc/smf/logger"
)


func ReadFromService(SUPI string) ([]byte, error) {
	os.Setenv("GODEBUG", "http2client=0")
	req, err1 := http.NewRequest("GET", "http://172.15.0.211:30000/links", nil)
	if err1 != nil {
		logger.CtxLog.Errorf(err1.Error())
		return nil, errors.New(err1.Error())
	}

	res, err2 := new(http.Client).Do(req)
	if err2 != nil {
		logger.CtxLog.Errorf(err2.Error())
		return nil, errors.New(err2.Error())
	}
	if res.StatusCode != http.StatusOK {
		logger.CtxLog.Errorf("None OK status code: %s", res.StatusCode)
		return nil, fmt.Errorf("None OK status code: %s", res.StatusCode)
	}

	resData, err3 := ioutil.ReadAll(res.Body)
	if err3 != nil {
		logger.CtxLog.Errorf(err3.Error())
		return nil, errors.New(err3.Error())
	}

	logger.CtxLog.Infof(string(resData))
	return resData, nil
}

func ReadUERoutesFromService() ([]byte, error) {
	os.Setenv("GODEBUG", "http2client=0")
	req, err1 := http.NewRequest("GET", "http://172.15.0.211:30000/ue-routes", nil)
	if err1 != nil {
		logger.CtxLog.Errorf(err1.Error())
		return nil, errors.New(err1.Error())
	}

	res, err2 := new(http.Client).Do(req)
	if err2 != nil {
		logger.CtxLog.Errorf(err2.Error())
		return nil, errors.New(err2.Error())
	}
	if res.StatusCode != http.StatusOK {
		logger.CtxLog.Errorf("None OK status code: %s", res.StatusCode)
		return nil, fmt.Errorf("None OK status code: %s", res.StatusCode)
	}

	resData, err3 := ioutil.ReadAll(res.Body)
	if err3 != nil {
		logger.CtxLog.Errorf(err3.Error())
		return nil, errors.New(err3.Error())
	}

	logger.CtxLog.Infof(string(resData))
	return resData, nil
}

func DynamicLoadLinksGET(SUPI string) error {
	logger.CtxLog.Infof("DynamicLoadLinksGET: SUPI = [%s]", SUPI)
	resData, err := ReadFromService(SUPI)
	if err != nil {
		return errors.New(err.Error())
	}
	// We fill into upi but only 'links' are currently applicable
	SmfLinksConf := factory.UserPlaneInformation{}
	if jsonErr := json.Unmarshal([]byte(resData), &SmfLinksConf); jsonErr != nil {
		// its from string type already
		return jsonErr
	}

	// Fill default topology into global SMF context
	ReloadLinks(SMF_Self().UserPlaneInformation, &SmfLinksConf)
	return nil
}

func DynamicLoadUERoutesGET() error {
	logger.CtxLog.Infof("DynamicLoadUERoutesGET")
	resData, err := ReadUERoutesFromService()
	if err != nil {
		return errors.New(err.Error())
	}
	SmfRoutingConf := factory.RoutingConfig{}
	if jsonErr := json.Unmarshal([]byte(resData), &SmfRoutingConf); jsonErr != nil {
		// its from string type already
		return jsonErr
	}

	// Fill default topology into global SMF context
	InitSMFUERouting(&SmfRoutingConf)
	return nil
}
