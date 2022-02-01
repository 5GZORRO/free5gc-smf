package upi

import (
	"net/http"
	"github.com/free5gc/smf/logger"
	"github.com/gin-gonic/gin"
)

type UserPlaneInformation struct {
	UPNodes map[string]UPNode `json:"up_nodes"`
}

// Should conform factory.config.UPNode
type UPNode struct {
	Type                 string                 `json:"type"`
	NodeID               string                 `json:"node_id"`
	ANIP                 string                 `json:"an_ip"`
	//SNssaiInfos          []SnssaiUpfInfoItem    `json:"sNssaiUpfInfos"`
	InterfaceUpfInfoList []InterfaceUpfInfoItem `json:"interfaces"`

//	User     string `form:"user" json:"user" binding:"required"`
//	Password string `form:"password" json:"password" binding:"required"`
}

type InterfaceUpfInfoItem struct {
//	InterfaceType   models.UpInterfaceType `json:"interfaceType"`
	Endpoints       []string               `json:"endpoints"`
	NetworkInstance string                 `json:"networkInstance"`
}

//type SnssaiUpfInfoItem struct {
//	SNssai         *models.Snssai   `json:"sNssai"`
//	DnnUpfInfoList []DnnUpfInfoItem `json:"dnnUpfInfoList"`
//}
//
//type DnnUpfInfoItem struct {
//	Dnn             string                  `json:"dnn"`
//	DnaiList        []string                `json:"dnaiList"`
//	PduSessionTypes []models.PduSessionType `json:"pduSessionTypes"`
//	Pools           []UEIPPool              `json:"pools"`
//}

func PostUpiUPFs(c *gin.Context) {
	logger.PduSessLog.Info("Recieve Add UPFs Request")
	var json UserPlaneInformation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if json.UPNodes["UPF-T5"].InterfaceUpfInfoList[1].Endpoints[0] != "172.15.0.29" || json.UPNodes["UPF-T5"].NodeID != "upf-t5.free5gc.org" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}