package context

import (
	"gopkg.in/yaml.v2"
	"github.com/free5gc/smf/factory"
)


const SmfConfig = `links:
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

func DynamicLoadLinks() error {
	// we use upi but only links applicable here....
	SmfLinksConf := factory.UserPlaneInformation{}
	if yamlErr := yaml.Unmarshal([]byte(SmfConfig), &SmfLinksConf); yamlErr != nil {
		return yamlErr
	}

	ReloadLinks(SMF_Self().UserPlaneInformation, &SmfLinksConf)
	return nil
}
