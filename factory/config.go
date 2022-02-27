/*
 * SMF Configuration Factory
 */

package factory

import (
	"time"

	"github.com/free5gc/logger_util"
	"github.com/free5gc/openapi/models"
)

const (
	SMF_EXPECTED_CONFIG_VERSION        = "1.0.2"
	UE_ROUTING_EXPECTED_CONFIG_VERSION = "1.0.1"
)

type Config struct {
	Info          *Info               `yaml:"info"`
	Configuration *Configuration      `yaml:"configuration"`
	Logger        *logger_util.Logger `yaml:"logger"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

const (
	SMF_DEFAULT_IPV4     = "127.0.0.2"
	SMF_DEFAULT_PORT     = "8000"
	SMF_DEFAULT_PORT_INT = 8000
)

type Configuration struct {
	SmfName              string               `yaml:"smfName,omitempty"`
	Sbi                  *Sbi                 `yaml:"sbi,omitempty"`
	PFCP                 *PFCP                `yaml:"pfcp,omitempty"`
	NrfUri               string               `yaml:"nrfUri,omitempty"`
	UserPlaneInformation UserPlaneInformation `yaml:"userplane_information"`
	ServiceNameList      []string             `yaml:"serviceNameList,omitempty"`
	SNssaiInfo           []SnssaiInfoItem     `yaml:"snssaiInfos,omitempty"`
	ULCL                 bool                 `yaml:"ulcl,omitempty"`
	PLMNList             []PLMNID             `yaml:"plmnList,omitempty"`
	Locality             string               `yaml:"locality,omitempty"`
}

type SnssaiInfoItem struct {
	SNssai   *models.Snssai      `yaml:"sNssai"`
	DnnInfos []SnssaiDnnInfoItem `yaml:"dnnInfos"`
}

type SnssaiDnnInfoItem struct {
	Dnn   string `yaml:"dnn"`
	DNS   *DNS   `yaml:"dns"`
	PCSCF *PCSCF `yaml:"pcscf,omitempty"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	TLS          *TLS   `yaml:"tls"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is registered at NRF.
	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"` // IP used to run the server in the node.
	Port        int    `yaml:"port,omitempty"`
}

type TLS struct {
	PEM string `yaml:"pem,omitempty"`
	Key string `yaml:"key,omitempty"`
}

type PFCP struct {
	Addr string `yaml:"addr,omitempty"`
	Port uint16 `yaml:"port,omitempty"`
}

type DNS struct {
	IPv4Addr string `yaml:"ipv4,omitempty"`
	IPv6Addr string `yaml:"ipv6,omitempty"`
}

type PCSCF struct {
	IPv4Addr string `yaml:"ipv4,omitempty" valid:"ipv4,required"`
}

type Path struct {
	DestinationIP   string   `yaml:"DestinationIP,omitempty"`
	DestinationPort string   `yaml:"DestinationPort,omitempty"`
	UPF             []string `yaml:"UPF,omitempty"`
}

type UERoutingInfo struct {
	Members       []string       `yaml:"members"`
	AN            string         `yaml:"AN,omitempty"`
	Topology      []UPLink       `yaml:"topology"`
	SpecificPaths []SpecificPath `yaml:"specificPath,omitempty"`
}

// RouteProfID is string providing a Route Profile identifier.
type RouteProfID string

// RouteProfile maintains the mapping between RouteProfileID and ForwardingPolicyID of UPF
type RouteProfile struct {
	// Forwarding Policy ID of the route profile
	ForwardingPolicyID string `yaml:"forwardingPolicyID,omitempty"`
}

// PfdContent represents the flow of the application
type PfdContent struct {
	// Identifies a PFD of an application identifier.
	PfdID string `yaml:"pfdID,omitempty"`
	// Represents a 3-tuple with protocol, server ip and server port for
	// UL/DL application traffic.
	FlowDescriptions []string `yaml:"flowDescriptions,omitempty"`
	// Indicates a URL or a regular expression which is used to match the
	// significant parts of the URL.
	Urls []string `yaml:"urls,omitempty"`
	// Indicates an FQDN or a regular expression as a domain name matching
	// criteria.
	DomainNames []string `yaml:"domainNames,omitempty"`
}

// PfdDataForApp represents the PFDs for an application identifier
type PfdDataForApp struct {
	// Identifier of an application.
	AppID string `yaml:"applicationId"`
	// PFDs for the application identifier.
	Pfds []PfdContent `yaml:"pfds"`
	// Caching time for an application identifier.
	CachingTime *time.Time `yaml:"cachingTime,omitempty"`
}

type RoutingConfig struct {
	Info          *Info                        `json:"info,omitempty" yaml:"info"`
	UERoutingInfo map[string]UERoutingInfo     `json:"ueRoutingInfo" yaml:"ueRoutingInfo"`
	RouteProf     map[RouteProfID]RouteProfile `json:"routeProfile,omitempty" yaml:"routeProfile,omitempty"`
	PfdDatas      []*PfdDataForApp             `json:"pfdDataForApp,omitempty" yaml:"pfdDataForApp,omitempty"`
}

// UserPlaneInformation describe core network userplane information
type UserPlaneInformation struct {
	UPNodes map[string]UPNode `json:"up_nodes,omitempty" yaml:"up_nodes,omitempty"`
	Links   []UPLink          `json:"links,omitempty" yaml:"links"`
}

// UPNode represent the user plane node
type UPNode struct {
	Type                 string                 `json:"type" yaml:"type"`
	NodeID               string                 `json:"node_id" yaml:"node_id"`
	ANIP                 string                 `json:"an_ip" yaml:"an_ip"`
	Dnn                  string                 `json:"dnn" yaml:"dnn"`
	SNssaiInfos          []SnssaiUpfInfoItem    `json:"sNssaiUpfInfos,omitempty" yaml:"sNssaiUpfInfos,omitempty"`
	InterfaceUpfInfoList []InterfaceUpfInfoItem `json:"interfaces,omitempty" yaml:"interfaces,omitempty"`
}

type InterfaceUpfInfoItem struct {
	InterfaceType   models.UpInterfaceType `json:"interfaceType" yaml:"interfaceType"`
	Endpoints       []string               `json:"endpoints" yaml:"endpoints"`
	NetworkInstance string                 `json:"networkInstance" yaml:"networkInstance"`
}

type SnssaiUpfInfoItem struct {
	SNssai         *models.Snssai   `json:"sNssai" yaml:"sNssai"`
	DnnUpfInfoList []DnnUpfInfoItem `json:"dnnUpfInfoList" yaml:"dnnUpfInfoList"`
}

type DnnUpfInfoItem struct {
	Dnn             string                  `json:"dnn" yaml:"dnn"`
	DnaiList        []string                `json:"dnaiList" yaml:"dnaiList"`
	PduSessionTypes []models.PduSessionType `json:"pduSessionTypes" yaml:"pduSessionTypes"`
	Pools           []UEIPPool              `json:"pools" yaml:"pools"`
}

type UPLink struct {
	A string `json:"A" yaml:"A"`
	B string `json:"B" yaml:"B"`
}

type UEIPPool struct {
	Cidr string `yaml:"cidr"`
}

type SpecificPath struct {
	DestinationIP   string   `yaml:"dest,omitempty"`
	DestinationPort string   `yaml:"DestinationPort,omitempty"`
	Path            []string `yaml:"path"`
}

type PLMNID struct {
	MCC string `yaml:"mcc"`
	MNC string `yaml:"mnc"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (r *RoutingConfig) GetVersion() string {
	if r.Info != nil && r.Info.Version != "" {
		return r.Info.Version
	}
	return ""
}

