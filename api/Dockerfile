# 1. sudo docker build --tag smf-api --force-rm=true .
# 2. sudo docker run -p 30000:8080 smf-api

FROM python:3.6.5-alpine

RUN apk add --no-cache bash \
 && apk add --no-cache --virtual .build-deps \
        bzip2-dev \
        g++ \
        libc-dev

RUN pip install --upgrade pip setuptools

RUN pip install --no-cache-dir gevent==1.2.1 flask>=1.0.0 requests==2.13.0 pyyaml \
  && apk del .build-deps

RUN mkdir -p /api
ADD server.py /api/

CMD ["/bin/bash", "-c", "cd api && python -u server.py"]
