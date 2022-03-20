# UPI

This is the UPI endpoint collection.

## API

### Add UPF

Add UPF definition to the datamodel

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

TODO

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

TODO
