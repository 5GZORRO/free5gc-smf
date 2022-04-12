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
	finalUri := fmt.Sprintf("%s/%s", SMF_Self().SmfExtUri, "links")
	req, err1 := http.NewRequest("GET", finalUri, nil)
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

func ReadUERoutesFromService(Uplink string, Downlink string) ([]byte, error) {
	logger.CtxLog.Infof("ReadUERoutesFromService: Uplink: %s, Downlink: %s", Uplink, Downlink)
	os.Setenv("GODEBUG", "http2client=0")
	finalUri := fmt.Sprintf("%s/%s", SMF_Self().SmfExtUri, "ue-routes")
	req, err1 := http.NewRequest("GET", finalUri, nil)
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

func DynamicLoadLinksGET(smContext *SMContext) error {
	logger.CtxLog.Infof("DynamicLoadLinksGET: SUPI = [%s]", smContext.Supi)
	resData, err := ReadFromService(smContext.Supi)
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

func DynamicLoadUERoutesGET(smContext *SMContext) error {
	logger.CtxLog.Infof("DynamicLoadUERoutesGET")
	resData, err := ReadUERoutesFromService(smContext.DnnConfiguration.SessionAmbr.Uplink,
		smContext.DnnConfiguration.SessionAmbr.Downlink)
	if err != nil {
		return errors.New(err.Error())
	}
	SmfRoutingConf := factory.RoutingConfig{}
	if jsonErr := json.Unmarshal([]byte(resData), &SmfRoutingConf); jsonErr != nil {
		// its from string type already
		return jsonErr
	}

	// Fill default topology into global SMF context
	InitSMFUERoutingInContext(smContext, &SmfRoutingConf)
	return nil
}

func InitSMFUERoutingInContext(smContext *SMContext, routingConfig *factory.RoutingConfig){

	UERoutingInfo := routingConfig.UERoutingInfo
	smContext.UEPreConfigPathPool = make(map[string]*UEPreConfigPaths)
	smContext.UEDefaultPathPool = make(map[string]*UEDefaultPaths)
	smContext.ULCLGroups = make(map[string][]string)

	for groupName, routingInfo := range UERoutingInfo {
		logger.CtxLog.Debugln("Set context for ULCL group: ", groupName)
		smContext.ULCLGroups[groupName] = routingInfo.Members
		uePreConfigPaths, err := NewUEPreConfigPaths(routingInfo.SpecificPaths)
		if err != nil {
			logger.CtxLog.Warnln(err)
		} else {
			smContext.UEPreConfigPathPool[groupName] = uePreConfigPaths
		}
		ueDefaultPaths, err := NewUEDefaultPaths(smfContext.UserPlaneInformation, routingInfo.Topology)
		if err != nil {
			logger.CtxLog.Warnln(err)
		} else {
			smContext.UEDefaultPathPool[groupName] = ueDefaultPaths
		}
	}
}
