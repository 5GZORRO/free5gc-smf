package context

import (
	"errors"
	"gopkg.in/yaml.v2"
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
	logger.CtxLog.Infof("DynamicLoadLinksGET: resData = [%s]", string(resData))
	// We fill into upi but only 'links' are currently applicable
	SmfLinksConf := factory.UserPlaneInformation{}
	if yamlErr := yaml.Unmarshal([]byte(resData), &SmfLinksConf); yamlErr != nil {
		// its from string type already
		return yamlErr
	}

	// Fill default topology into global SMF context
	ReloadLinks(SMF_Self().UserPlaneInformation, &SmfLinksConf)
	return nil
}
