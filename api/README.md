# SMF-ext

This is an extension to SMF. The purpose of this extension is to provide specific data paths to be used for
session creation.

## Deploy the service

The service is jointly deployed with smf. Refer to free5gc [docker-compose](https://github.ibm.com/WEIT/free5gc-compose/tree/e0d4742-dynamic_load) setup

## API - default topology

### Update default topology

Create or update default topology

```
curl -H "Content-type: application/json" -X POST http://smf_api_address:8080/links
```

REST path:

```
    smf_api_ip_address - ipaddress of SMF-ext service
```

Return:

```
    status - 200
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X POST -d "@payloads/links.json" http://127.0.0.1:8080/links
```

### Get default topology

Get default topology

```
curl -H "Content-type: application/json" -X GET http://smf_api_address:8080/links
```

REST path:

```
    smf_api_ip_address - ipaddress of SMF-ext service
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X GET http://127.0.0.1:8080/links

    {
      "links": [
        {
          "A": "gNB1",
          "B": "UPF-R1"
        },
        {
          "A": "UPF-R1",
          "B": "UPF-T1"
        },
        {
          "A": "UPF-T1",
          "B": "UPF-C1"
        }
      ]
    }
```

## API - specific topology


### Create topology group

Create a group

```
curl -H "Content-type: application/json" -X POST http://smf_api_address:8080/ue-routes/<group_name>
```

REST path:

```
    smf_api_ip_address - ipaddress of SMF-ext service
    group_name         - the name of the group to add (str)
```

Return:

```
    status - 200
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X POST http://smf_api_address:8080/ue-routes/red
```

### Get group info

Returns information for the given group

```
curl -H "Content-type: application/json" -GET http://smf_api_address:8080/ue-routes/<group_name>
```

REST path:

```
    smf_api_ip_address - ipaddress of SMF-ext service
    group_name         - the name of the gruop (str)
```

Return:

```
    status - 200
    group information including members and topology (dict)
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X GET http://127.0.0.1:8080/ue-routes/red

    {
      "members": [
        "imsi-208938080000001"
      ],
      "specificPath": [
        {
          "dest": "60.61.0.0/16",
          "path": [
            "UPF-R1",
            "UPF-T1",
            "UPF-C1"
          ]
        }
      ],
      "topology": [
        {
          "A": "gNB1",
          "B": "UPF-R1"
        },
        {
          "A": "UPF-R1",
          "B": "UPF-T1"
        },
        {
          "A": "UPF-T1",
          "B": "UPF-C1"
        }
      ]
    }
```


### Add group member

Add a member to a given group

```
curl -H "Content-type: application/json" -X POST http://smf_api_address:8080/ue-routes/<group_name>/members/<member_name>
```

REST path:

```
    smf_api_ip_address - ipaddress of SMF-ext service
    group_name         - the name of the group (str)
    member_name        - the name of the member to add (str)
```

Return:

```
    status - 200
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X POST http://127.0.0.1:8080/ue-routes/red/members/imsi-208938080000007
```


### List group members

Returns members for the given group

```
curl -H "Content-type: application/json" -GET http://smf_api_address:8080/ue-routes/<group_name>/members
```

REST path:

```
    smf_api_ip_address - ipaddress of SMF-ext service
    group_name         - the name of the gruop (str)
```

Return:

```
    status - 200
    list of members (str)
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X GET http://127.0.0.1:8080/ue-routes/red/members

    [
      "imsi-208938080000001",
      "imsi-208938080000004",
      "imsi-208938080000007",
      "imsi-208938080000008"
    ]
```


### Add topology

Add topology for a given group

**Important:** it is mandatory for the UPFs to be pre-registered in SMF (see: [here](https://github.ibm.com/WEIT/free5gc-compose/blob/e0d4742-dynamic_load/config/smf/README.md))

```
curl -H "Content-type: application/json" -X POST http://smf_api_address:8080/ue-routes/<group_name>/topology
```

REST path:

```
    smf_api_ip_address - ipaddress of SMF-ext service
    group_name         - the name of the group (str)
```

Return:

```
    status - 200
```

Invocation example:

```bash
curl -X POST \
  http://127.0.0.1:8080/ue-routes/red/topology \
  -H 'content-type: application/json' \
  -d '{
  "topology": [
    {"A": "gNB1", "B": "UPF-R1"},
    {"A": "UPF-R1", "B": "UPF-T1"},
    {"A": "UPF-T1", "B": "UPF-C1"}
  ],
  "specificPath": [
    {"dest": "60.61.0.0/16", "path": ["UPF-R1", "UPF-T1", "UPF-C1"]}
  ]
}'
```

Alternatively - pass json file

```bash
    curl -H "Content-type: application/json" -X POST -d "@payloads/red-topology.json" http://127.0.0.1:8080/ue-routes/red/topology
```


### Update topology with link

Add a single link to a given topology

**Note:** group entry is created if does not exist

```
curl -H "Content-type: application/json" -X PUT http://smf_api_address:8080/ue-routes/<group_name>/topology
```

REST path:

```
    smf_api_ip_address - ipaddress of SMF-ext service
    group_name         - the name of the group (str)
```

Return:

```
    status - 200
```

Invocation examples:

define path: gNB1 -> UPF-R1 -> UPF-T1 -> UPF-C2

```bash
curl -X PUT \
  http://127.0.0.1:8080/ue-routes/red/topology \
  -H 'content-type: application/json' \
  -d '{"A": "gNB1", "B": "UPF-R1"}}'
```


```bash
curl -X PUT \
  http://127.0.0.1:8080/ue-routes/red/topology \
  -H 'content-type: application/json' \
  -d '{"A": "UPF-R1", "B": "UPF-T1"}}'
```

```bash
curl -X PUT \
  http://127.0.0.1:8080/ue-routes/red/topology \
  -H 'content-type: application/json' \
  -d '{"A": "UPF-T1", "B": "UPF-C2"}}'
```
