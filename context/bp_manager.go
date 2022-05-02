package context

import (
	"github.com/free5gc/smf/logger"
	"reflect"
)

type BPManager struct {
	BPStatus       BPStatus
	AddingPSAState AddingPSAState
	// Need these variable conducting Add additional PSA (TS23.502 4.3.5.4)
	// There value will change from time to time

	PendingUPF            PendingUPF
	ActivatedPaths        []*DataPath
	ActivatingPath        *DataPath
	UpdatedBranchingPoint map[*UPF]int
	ULCL                  *UPF
}
type BPStatus int

const (
	UnInitialized BPStatus = iota
	AddingPSA
	AddPSASuccess
	InitializedSuccess
	InitializedFail
)

type AddingPSAState int

const (
	ActivatingDataPath AddingPSAState = iota
	EstablishingNewPSA
	EstablishingULCL
	UpdatingPSA2DownLink
	UpdatingRANAndIUPFUpLink
	Finished
)

type PendingUPF map[string]bool

func NewBPManager(supi string) (bpManager *BPManager) {
	bpManager = &BPManager{
		BPStatus:              UnInitialized,
		AddingPSAState:        ActivatingDataPath,
		ActivatedPaths:        make([]*DataPath, 0),
		UpdatedBranchingPoint: make(map[*UPF]int),
		PendingUPF:            make(PendingUPF),
	}

	return
}

func (bpMGR *BPManager) SelectPSA2(smContext *SMContext) {
	hasSelectPSA2 := false
	bpMGR.ActivatedPaths = []*DataPath{}
	// smContext.Tunnel.DataPathPool: first entry should point to the default
	// path derived from the ue-route topology
	// Next ones are the preConfig paths
	for _, dataPath := range smContext.Tunnel.DataPathPool {
		if dataPath.Activated {
			logger.PfcpLog.Traceln("SelectPSA2: Add to ActivatedPaths:\n" + dataPath.String() + "\n")
			bpMGR.ActivatedPaths = append(bpMGR.ActivatedPaths, dataPath)
		} else {
			if !hasSelectPSA2 {
				bpMGR.ActivatingPath = dataPath
				logger.PfcpLog.Traceln("SelectPSA2: Add to ActivatingPath:\n" + dataPath.String() + "\n")
				// It seems to select a preConfig path
				hasSelectPSA2 = true
			}
		}
	}
}

func (bpMGR *BPManager) FindULCL(smContext *SMContext) error {
	bpMGR.UpdatedBranchingPoint = make(map[*UPF]int)
	// this is the preConfig path
	activatingPath := bpMGR.ActivatingPath
	// this is the default selected path that already established per this session
	for _, psa1Path := range bpMGR.ActivatedPaths {
		depth := 0
		// the 1st upf in the established path
		psa1CurDPNode := psa1Path.FirstDPNode
		for psa2CurDPNode := activatingPath.FirstDPNode; psa2CurDPNode != nil; psa2CurDPNode = psa2CurDPNode.Next() {
			// if the node in preConfig path is the 1st upf in the established one
			if reflect.DeepEqual(psa2CurDPNode.UPF.NodeID, psa1CurDPNode.UPF.NodeID) {
				psa1CurDPNode = psa1CurDPNode.Next()
				depth++

				if _, exist := bpMGR.UpdatedBranchingPoint[psa2CurDPNode.UPF]; !exist {
					bpMGR.UpdatedBranchingPoint[psa2CurDPNode.UPF] = depth
				}
			} else {
				break
			}
		}
	}

	maxDepth := 0
	for upf, depth := range bpMGR.UpdatedBranchingPoint {
		if depth > maxDepth {
			bpMGR.ULCL = upf
			maxDepth = depth
		}
	}
	if bpMGR.ULCL != nil {
		logger.CtxLog.Warnf("bpMGR.ULCL:[%+v]", bpMGR.ULCL.NodeID.ResolveNodeIdToIp().To4())
	} else {
		logger.CtxLog.Warnf("bpMGR.ULCL NULL !!")
	}
	return nil
}

func (pendingUPF PendingUPF) IsEmpty() bool {
	if len(pendingUPF) == 0 {
		return true
	} else {
		return false
	}
}
