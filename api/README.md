# SMF-api

## Deploy the service

```
docker build --tag smf-api --force-rm=true .
docker run -p 30000:8080 smf-api
```

## API

### Update links

```
curl -H "Content-type: application/json" -X POST -d "@links1.json" http://172.15.0.211:30000/links
```

### Get links

```
curl -H "Content-type: application/json" -X GET http://172.15.0.211:30000/links
```
