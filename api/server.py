import flask
import json
import os
import requests
import sys
import yaml


from gevent.wsgi import WSGIServer
from werkzeug.exceptions import HTTPException

class Proxy:
    def __init__(self):
        sys.stdout.write('SMF-API initialized\n')


proxy = flask.Flask(__name__)
proxy.debug = True
server = None

proxy_server = None


links = dict()


def setServer(s):
    global server
    server = s


def setProxy(p):
    global proxy_server
    proxy_server = p


def getMessagePayload():
    message = flask.request.get_json(force=True, silent=True)
    if message and not isinstance(message, dict):
        flask.abort(400, 'message payload is not a dictionary')
    else:
        value = message if (message or message == {}) else {}
    if not isinstance(value, dict):
        flask.abort(400, 'message payload did not provide binding for "value"')
    return value


@proxy.route("/hello")
def hello():
    sys.stdout.write ('Enter /hello\n')
    return ("Greetings from the SMF-API server! ")


@proxy.route("/links", methods=['POST'])
def links_create():
    sys.stdout.write ('Enter POST /links\n')
    value = getMessagePayload()
    global links
    links = value
    return yaml.dump(links)
#     return """
# links:
#   - A: gNB1
#     B: UPF-R1
# 
#   - A: UPF-R1
#     B: UPF-C3
# """



@proxy.route("/links", methods=['GET'])
def links():
    sys.stdout.write ('Enter /links\n')
    return yaml.dump(links)
# """
# links:
#   - A: gNB1
#     B: UPF-R1
# 
#   - A: UPF-R1
#     B: UPF-C3
# """



def main():
    port = int(os.getenv('LISTEN_PORT', 8080))
    server = WSGIServer(('0.0.0.0', port), proxy, log=None)
    setServer(server)
    print ('\n\n-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-')
    print ("Starting SMF-API .. ready to serve requests on PORT: %s..\n\n" %
           (int(port)))
    print ('-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n\n')

    server.serve_forever()


if __name__ == '__main__':
    setProxy(Proxy())
    main()
