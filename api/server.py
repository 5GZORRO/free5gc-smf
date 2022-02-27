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


links = {}
ueroutes = {
    "ueRoutingInfo": {
        "east": {
            "members": [
                'imsi-208930000000007'
            ],
            "topology": [
                {"A": "gNB1", "B": "UPF-R1"},
                {"A": "UPF-R1", "B": "UPF-T1"},
                {"A": "UPF-T1", "B": "UPF-C1"}
            ]
        },
        "west": {
            "members": [
                'imsi-208930000000001',
                'imsi-208930000000003'
            ],
            "topology":[
                {"A": "gNB1", "B": "UPF-R1"},
                {"A": "UPF-R1", "B": "UPF-T2"},
                {"A": "UPF-T2", "B": "UPF-C4"}
            ]
        },
        "middle": {
            "members": [
                'imsi-208930000000004',
                'imsi-208930000000005',
                'imsi-208930000000002'
            ],
            "topology": [
                {"A": "gNB1", "B": "UPF-R1"},
                {"A": "UPF-R1", "B": "UPF-T2"},
                {"A": "UPF-T2", "B": "UPF-C2"}
            ]
        }
    }
}



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
def links_post():
    sys.stdout.write ('Enter POST /links\n')
    value = getMessagePayload()
    global links
    links = value
    return flask.jsonify(links)


@proxy.route("/links", methods=['GET'])
def links_get():
    sys.stdout.write ('Enter /links\n')
    global links
    if not links:
        response = flask.jsonify({'NOT_FOUND': 404})
        response.status_code = 404
        return response
    else:
        return flask.jsonify(links)


@proxy.route("/ue-routes", methods=['GET'])
def ueroutes_get():
    sys.stdout.write ('Enter /ur-routes\n')
    global ueroutes
    if not ueroutes:
        response = flask.jsonify({'NOT_FOUND': 404})
        response.status_code = 404
        return response
    else:
        return flask.jsonify(ueroutes)


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
