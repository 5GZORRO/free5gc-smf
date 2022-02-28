# SMF-api

## Deploy the service

```
docker build --tag smf-api --force-rm=true .
docker run -p 30000:8080 smf-api
```

## API

### Update default topology

Create or update default topology

```
curl -H "Content-type: application/json" -X POST http://smf_api_address:30000/links
```

REST path:

```
    smf_api_ip_address - ipaddress SMF API service
```

Return:

```
    status - 200
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X POST -d "@payloads/links.json" http://172.15.0.211:30000/links
```

### Get default topology

Get default topology

```
curl -H "Content-type: application/json" -X GET http://smf_api_address:30000/links
```

REST path:

```
    smf_api_ip_address - ipaddress SMF API service
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X GET http://172.15.0.211:30000/links

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


### Create group

Create a group

```
curl -H "Content-type: application/json" -X POST http://smf_api_address:30000/ue-routes/<group_name>
```

REST path:

```
    smf_api_ip_address - ipaddress SMF API service
    group_name         - the name of the group to add (str)
```

Return:

```
    status - 200
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X POST http://smf_api_address:30000/ue-routes/east
```

### Get group info

Returns information for the given group

```
curl -H "Content-type: application/json" -GET http://smf_api_address:30000/ue-routes/<group_name>
```

REST path:

```
    smf_api_ip_address - ipaddress SMF API service
    group_name         - the name of the gruop (str)
```

Return:

```
    status - 200
    group information including members and topology (dict)
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X GET http://172.15.0.211:30000/ue-routes/east

    {
      "members": [
        "imsi-208930000000001"
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
curl -H "Content-type: application/json" -X POST http://smf_api_address:30000/ue-routes/<group_name>/members/<member_name>
```

REST path:

```
    smf_api_ip_address - ipaddress SMF API service
    group_name         - the name of the group (str)
    member_name        - the name of the member to add (str)
```

Return:

```
    status - 200
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X POST http://172.15.0.211:30000/ue-routes/east/members/imsi-208930000000007
```


### List group members

Returns members for the given group

```
curl -H "Content-type: application/json" -GET http://smf_api_address:30000/ue-routes/<group_name>/members
```

REST path:

```
    smf_api_ip_address - ipaddress SMF API service
    group_name         - the name of the gruop (str)
```

Return:

```
    status - 200
    list of members (str)
```

Invocation example:

```bash
    curl -H "Content-type: application/json" -X GET http://172.15.0.211:30000/ue-routes/east/members

    [
      "imsi-208930000000001",
      "imsi-208930000000004",
      "imsi-208930000000007",
      "imsi-208930000000008"
    ]
```


### Add topology

Add topology for a given group

**Important:** it is mandatory for the UPFs to be pre-registered in SMF (see: [here](https://github.ibm.com/WEIT/free5gc-compose/blob/e0d4742-dynamic_load/config/smf/README.md))

```
curl -H "Content-type: application/json" -X POST http://smf_api_address:30000/ue-routes/<group_name>/topology
```

REST path:

```
    smf_api_ip_address - ipaddress SMF API service
    group_name         - the name of the group (str)
```

Return:

```
    status - 200
```

Invocation example:

```bash
curl -X POST \
  http://172.15.0.211:30000/ue-routes/east/topology \
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
    curl -H "Content-type: application/json" -X POST -d "@payloads/east-topology.json" http://172.15.0.211:30000/ue-routes/east/topology
```
