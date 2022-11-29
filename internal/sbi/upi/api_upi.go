package upi

import (
    "net/http"
    "net"
    "github.com/free5gc/smf/internal/logger"
    "github.com/gin-gonic/gin"
    smf_context "github.com/free5gc/smf/internal/context"
    "github.com/free5gc/pfcp/pfcpType"
    "github.com/free5gc/smf/internal/pfcp/message"
    "github.com/free5gc/smf/pkg/factory"
    "github.com/free5gc/openapi/models"
    "github.com/free5gc/util/httpwrapper"
)

func DeleteUpiUpf(c *gin.Context) {
    req := httpwrapper.NewRequest(c.Request, nil)
    req.Params["upfRef"] = c.Params.ByName("upfRef")

    upfRef := req.Params["upfRef"]
    upi := smf_context.SMF_Self().UserPlaneInformation
    found := false

    upNode, ok := upi.UPNodes[upfRef]
    if ok {
        found = true
        logger.InitLog.Infof("UPF [%s] FOUND and is about to get removed.\n", upfRef)
        if upNode.Type == "UPF" {
            smf_context.RemoveUPFNodeByNodeID(upNode.UPF.NodeID)
        }
        delete(upi.UPNodes, upfRef)
    }

    if found {
        c.JSON(http.StatusNoContent, gin.H{})
    } else {
        c.JSON(http.StatusNotFound, gin.H{})
    }
}


func GetUpi(c *gin.Context) {
        upi := smf_context.SMF_Self().UserPlaneInformation
        nodes := make(map[string]factory.UPNode)
        for name, upNode := range upi.UPNodes {
                u := new(factory.UPNode)
                switch upNode.Type {
                    case smf_context.UPNODE_UPF:
                    u.Type = "UPF"
                    case smf_context.UPNODE_AN:
                    u.Type = "AN"
                    u.ANIP = upNode.ANIP.String()
                    u.NrCellId = upNode.NrCellId
                    default:
                    u.Type = "Unkown"
                }

                nodeIDtoIp := upNode.NodeID.ResolveNodeIdToIp()
                // for AN nodeIDtoIp is nil
                if nodeIDtoIp != nil {
                    u.NodeID = nodeIDtoIp.String()
                }
                // for AN UPF is nil
                if upNode.UPF != nil {
                    if upNode.UPF.SNssaiInfos != nil {
                        FsNssaiInfoList := make([]factory.SnssaiUpfInfoItem, 0)
                        for _, sNssaiInfo := range upNode.UPF.SNssaiInfos {
                            FDnnUpfInfoList := make([]factory.DnnUpfInfoItem, 0)
                            for _, dnnInfo := range sNssaiInfo.DnnList {
                                FUEIPPools := make([]factory.UEIPPool, 0)
                                for _, pool := range dnnInfo.UeIPPools {
                                    FUEIPPools = append(FUEIPPools, factory.UEIPPool{
                                            Cidr: pool.UeSubNet.String(),
                                        })
                                } // for pool
                                FDnnUpfInfoList = append(FDnnUpfInfoList, factory.DnnUpfInfoItem{
                                    Dnn: dnnInfo.Dnn,
                                    Pools: FUEIPPools,
                                })
                            } // for dnnInfo
                            Fsnssai := factory.SnssaiUpfInfoItem {
                                SNssai: &models.Snssai{
                                    Sst: sNssaiInfo.SNssai.Sst,
                                    Sd: sNssaiInfo.SNssai.Sd,
                                },
                                DnnUpfInfoList: FDnnUpfInfoList,
                             }
                            FsNssaiInfoList = append(FsNssaiInfoList, Fsnssai)
                        } // for sNssaiInfo
                        u.SNssaiInfos = FsNssaiInfoList
                    } // if UPF.SNssaiInfos

                    FNxList := make([]factory.InterfaceUpfInfoItem, 0)
                    for _, iface := range upNode.UPF.N3Interfaces {
                        endpoints := make([]string, 0)
                        // upf.go L90
                        if iface.EndpointFQDN != "" {
                            endpoints = append(endpoints, iface.EndpointFQDN)
                        }
                        for _, eIP := range iface.IPv4EndPointAddresses {
                            endpoints = append(endpoints, eIP.String())
                        }
                        FNxList = append(FNxList, factory.InterfaceUpfInfoItem{
                            InterfaceType: models.UpInterfaceType_N3,
                            Endpoints: endpoints,
                            NetworkInstance: iface.NetworkInstance,
                        })
                    } // for N3Interfaces

                    //FN9List := make([]factory.InterfaceUpfInfoItem, 0)
                    for _, iface := range upNode.UPF.N9Interfaces {
                        endpoints := make([]string, 0)
                        // upf.go L90
                        if iface.EndpointFQDN != "" {
                            endpoints = append(endpoints, iface.EndpointFQDN)
                        }
                        for _, eIP := range iface.IPv4EndPointAddresses {
                            endpoints = append(endpoints, eIP.String())
                        }
                        FNxList = append(FNxList, factory.InterfaceUpfInfoItem{
                            InterfaceType: models.UpInterfaceType_N9,
                            Endpoints: endpoints,
                            NetworkInstance: iface.NetworkInstance,
                        })
                    } // N9Interfaces
                    u.InterfaceUpfInfoList = FNxList
                }
                nodes[name] = *u
        }

        json := &factory.UserPlaneInformation{
            UPNodes: nodes,
        }

        httpResponse := &httpwrapper.Response{
                Header: nil,
                Status: http.StatusOK,
                Body: json,
        }

        c.JSON(httpResponse.Status, httpResponse.Body)
}

func AddUPFs(upi *smf_context.UserPlaneInformation, upTopology *factory.UserPlaneInformation) {
    for name, node := range upTopology.UPNodes {
        if _, ok := upi.UPNodes[name]; ok {
            // TODO: consider it as an error?
            logger.InitLog.Warningf("UPF [%s] already exists in SMF. Ignoring request.\n", name)
            continue
        }
        upNode := new(smf_context.UPNode)
        upNode.Type = smf_context.UPNodeType(node.Type)
        switch upNode.Type {
        case smf_context.UPNODE_UPF:
            // ParseIp() always return 16 bytes
            // so we can't use the length of return ip to separate IPv4 and IPv6
            // This is just a work around
            var ip net.IP
            if net.ParseIP(node.NodeID).To4() == nil {
                ip = net.ParseIP(node.NodeID)
            } else {
                ip = net.ParseIP(node.NodeID).To4()
            }

            switch len(ip) {
            case net.IPv4len:
                upNode.NodeID = pfcpType.NodeID{
                    NodeIdType:  pfcpType.NodeIdTypeIpv4Address,
                    IP: ip,
                }
            case net.IPv6len:
                upNode.NodeID = pfcpType.NodeID{
                    NodeIdType:  pfcpType.NodeIdTypeIpv6Address,
                    IP: ip,
                }
            default:
                upNode.NodeID = pfcpType.NodeID{
                    NodeIdType:  pfcpType.NodeIdTypeFqdn,
                    FQDN:       node.NodeID,
                }
            }

            upNode.UPF = smf_context.NewUPF(&upNode.NodeID, node.InterfaceUpfInfoList)
            snssaiInfos := make([]smf_context.SnssaiUPFInfo, 0)
            for _, snssaiInfoConfig := range node.SNssaiInfos {
                snssaiInfo := smf_context.SnssaiUPFInfo{
                    SNssai: smf_context.SNssai{
                        Sst: snssaiInfoConfig.SNssai.Sst,
                        Sd:  snssaiInfoConfig.SNssai.Sd,
                    },
                    DnnList: make([]smf_context.DnnUPFInfoItem, 0),
                }

                for _, dnnInfoConfig := range snssaiInfoConfig.DnnUpfInfoList {
                    ueIPPools := make([]*smf_context.UeIPPool, 0)
                    for _, pool := range dnnInfoConfig.Pools {
                        ueIPPool := smf_context.NewUEIPPool(&pool)
                        if ueIPPool == nil {
                            logger.InitLog.Fatalf("invalid pools value: %+v", pool)
                        } else {
                            ueIPPools = append(ueIPPools, ueIPPool)
                            /* TODO: check overlapping cidrs*/
                            // allUEIPPools = append(allUEIPPools, ueIPPool)
                        }
                    }
                    snssaiInfo.DnnList = append(snssaiInfo.DnnList, smf_context.DnnUPFInfoItem{
                        Dnn:             dnnInfoConfig.Dnn,
                        DnaiList:        dnnInfoConfig.DnaiList,
                        PduSessionTypes: dnnInfoConfig.PduSessionTypes,
                        UeIPPools:       ueIPPools,
                    })
                }
                snssaiInfos = append(snssaiInfos, snssaiInfo)
            }
            upNode.UPF.SNssaiInfos = snssaiInfos
            upi.UPFs[name] = upNode
        case smf_context.UPNODE_AN:
            upNode.ANIP = net.ParseIP(node.ANIP)
            upNode.NrCellId = node.NrCellId
            upi.AccessNetwork[name] = upNode
        default:
            logger.InitLog.Warningf("invalid UPNodeType: %s\n", upNode.Type)
        }

        upi.UPNodes[name] = upNode

        ipStr := upNode.NodeID.ResolveNodeIdToIp().String()
        upi.UPFIPToName[ipStr] = name

        // AllocateUPFID
        upfid := upNode.UPF.UUID()
        upfip := upNode.NodeID.ResolveNodeIdToIp().String()
        upi.UPFsID[name] = upfid
        upi.UPFsIPtoID[upfip] = upfid

        // Association (asynch)
        // TODO: should it be here?
        upf := upNode.UPF
        if upf.NodeID.NodeIdType == pfcpType.NodeIdTypeFqdn {
            logger.AppLog.Infof("Send PFCP Association Request to UPF[%s](%s)\n", upf.NodeID.FQDN,
                upf.NodeID.ResolveNodeIdToIp().String())
        } else {
            logger.AppLog.Infof("Send PFCP Association Request to UPF[%s]\n", upf.NodeID.IP)
        }
        message.SendPfcpAssociationSetupRequest(upf.NodeID)
    }
}

func PostUpiUPFs(c *gin.Context) {
    logger.PduSessLog.Info("Recieve Add UPFs Request")
    var json factory.UserPlaneInformation
    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    logger.PduSessLog.Info("About to add UPFs")
    AddUPFs(smf_context.SMF_Self().UserPlaneInformation, &json)

    c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
