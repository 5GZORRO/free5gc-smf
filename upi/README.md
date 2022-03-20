# UPI

This is the UPI endpoint collection.

## API

### Add UPF

Add UPF definition to the datamodel.

**Note:** SMF must access UPF `node_id` (N4) so that it can register it

```
curl -H "Content-type: application/json" -X POST '{...}' http://smf_address:8000/upi/v1/upf
```

REST path:

```
    smf_ip_address - SMF ipaddress
```

Return:

```
    status - 200
```

Invocation example:

```bash
curl -X POST \
  http://127.0.0.2:8000/upi/v1/upf \
  -H "Content-type: application/json" \
-d '{
    "up_nodes": {
        "UPF-C1": {
            "type": "UPF",
            "node_id": "upf-c1.free5gc.org",
            "sNssaiUpfInfos": [
              {
                  "sNssai": {"sst": 1, "sd": "010203"},
                  "dnnUpfInfoList": [
                      {
                          "dnn": "internet",
                          "pools": [{"cidr": "60.61.0.0/24"}]
                      }
                  ]
              }
            ],
            "interfaces": [
                {
                    "interfaceType": "N9",
                    "endpoints": ["172.15.0.9"],
                    "networkInstance": "internet"
                }            
            ]
        }
    }
}'

{"status":"OK"}
```

### Dump UPFs

Get UPF datamodel

```
curl -H "Content-type: application/json" -X GET http://smf_address:8000/upi/v1/upf
```

REST path:

```
    smf_ip_address - SMF ipaddress
```

Invocation example:

```bash
curl http://127.0.0.2:8000/upi/v1/upf

{
  "up_nodes": {
    "UPF": {
      "type": "UPF",
      "node_id": "127.0.0.8",
      "an_ip": "",
      "dnn": "",
      "sNssaiUpfInfos": [
        {
          "sNssai": {
            "sst": 1,
            "sd": "010203"
          },
          "dnnUpfInfoList": [
            {
              "dnn": "internet",
              "dnaiList": null,
              "pduSessionTypes": null,
              "pools": [
                {
                  "Cidr": "60.60.0.0/16"
                }
              ]
            }
          ]
        },
        {
          "sNssai": {
            "sst": 1,
            "sd": "112233"
          },
          "dnnUpfInfoList": [
            {
              "dnn": "internet",
              "dnaiList": null,
              "pduSessionTypes": null,
              "pools": [
                {
                  "Cidr": "60.61.0.0/16"
                }
              ]
            }
          ]
        }
      ],
      "interfaces": [
        {
          "interfaceType": "N3",
          "endpoints": [
            "127.0.0.8"
          ],
          "networkInstance": "internet"
        }
      ]
    },
    "UPF-C1": {
      "type": "UPF",
      "node_id": "upf-c1.free5gc.org",
      "an_ip": "",
      "dnn": "",
      "sNssaiUpfInfos": [
        {
          "sNssai": {
            "sst": 1,
            "sd": "010203"
          },
          "dnnUpfInfoList": [
            {
              "dnn": "internet",
              "dnaiList": null,
              "pduSessionTypes": null,
              "pools": [
                {
                  "Cidr": "60.61.0.0/24"
                }
              ]
            }
          ]
        }
      ],
      "interfaces": [
        {
          "interfaceType": "N9",
          "endpoints": [
            "172.15.0.9"
          ],
          "networkInstance": "internet"
        }
      ]
    },
    "gNB1": {
      "type": "AN",
      "node_id": "",
      "an_ip": "<nil>",
      "dnn": ""
    }
  }
}

```

