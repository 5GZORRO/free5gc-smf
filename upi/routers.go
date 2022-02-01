package upi

import (
	"net/http"

	"strings"
	"github.com/gin-gonic/gin"

	"github.com/free5gc/logger_util"
	"github.com/free5gc/smf/logger"
)

//type UserPlaneInformationUPFsData struct {
//	UPNodes map[string]UPNode `json:"up_nodes"`
//}
//
//// UPNode represent the user plane node
//type UPNode struct {
//	Type                 string                 `json:"type"`
//	NodeID               string                 `json:"node_id"`
//	ANIP                 string                 `json:"an_ip"`
//	Dnn                  string                 `json:"dnn"`
//	SNssaiInfos          []SnssaiUpfInfoItem    `json:"sNssaiUpfInfos,omitempty"`
//	InterfaceUpfInfoList []InterfaceUpfInfoItem `json:"interfaces,omitempty"`
//}
//
//type SnssaiUpfInfoItem struct {
//	SNssai         *models.Snssai   `json:"sNssai"`
//	DnnUpfInfoList []DnnUpfInfoItem `json:"dnnUpfInfoList"`
//}
//
//type InterfaceUpfInfoItem struct {
//	InterfaceType   models.UpInterfaceType `json:"interfaceType"`
//	Endpoints       []string               `json:"endpoints"`
//	NetworkInstance string                 `json:"networkInstance"`
//}
//
//type DnnUpfInfoItem struct {
//	Dnn             string                  `json:"dnn"`
//	DnaiList        []string                `json:"dnaiList"`
//	PduSessionTypes []models.PduSessionType `json:"pduSessionTypes"`
//	Pools           []UEIPPool              `json:"pools"`
//}


// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)
	AddService(router)
	return router
}

func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("/upi/v1")

	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.HandlerFunc)
		}
	}
	return group
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/",
		Index,
	},

	{
		"PostUpiUPFs",
		strings.ToUpper("Post"),
		"/upf",
		PostUpiUPFs,
	},
}