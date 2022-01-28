package context

import (
	"gopkg.in/yaml.v2"
	"github.com/free5gc/smf/factory"
)


const SmfConfig = `userplane_information:
    up_nodes:
      gNB1:
        type: AN
        an_ip: 172.15.0.211

      UPF-R1:
        type: UPF
        node_id: upf-r1.free5gc.org
        sNssaiUpfInfos:
          - sNssai:
              sst: 1
              sd: 010203
            dnnUpfInfoList:
              - dnn: internet
        interfaces:
          - interfaceType: N3
            endpoints:
              - 172.15.0.6
            networkInstance: internet
          - interfaceType: N9
            endpoints:
              - 172.15.0.6
            networkInstance: internet

      UPF-C3:
        type: UPF
        node_id: upf-c3.free5gc.org
        sNssaiUpfInfos:
          - sNssai:
              sst: 1
              sd: 010203
            dnnUpfInfoList:
              - dnn: internet
                pools:
                  - cidr: 60.63.0.0/24
        interfaces:
          - interfaceType: N9
            endpoints:
              - 172.15.0.11
            networkInstance: internet
    links:
      - A: gNB1
        B: UPF-R1

      - A: UPF-R1
        B: UPF-C3
`

const UERoutingConfig = `ueRoutingInfo:
  UE-gNB-1-1:
    members:
    - imsi-208930000000002
    - imsi-208930000000001
    topology:
      - A: gNB1
        B: UPF-R1

      - A: UPF-R1
        B: UPF-C3

    specificPath:
      - dest: 60.61.0.0/16
        path: [UPF-R1, UPF-C3]
`

func DynamicLoad() error {
//        if _, err := template.New("blah").Parse(SmfConfig); err != nil {
//                return err
//        }

        SmfConfiguration := factory.Configuration{}

        if yamlErr := yaml.Unmarshal([]byte(SmfConfig), &SmfConfiguration); yamlErr != nil {
                return yamlErr
        }

        SMF_Self().UserPlaneInformation = NewUserPlaneInformation(&SmfConfiguration.UserPlaneInformation)
        return nil
}
