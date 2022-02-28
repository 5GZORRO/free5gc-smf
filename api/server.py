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


def find(l, predicate):
    results = [x for x in l if predicate(x)]
    return results[0] if len(results) > 0 else None


links = {}

prefix = {
    "info": {
        "version": "1.0.1",
        "description": "Routing information for UE"
    }
}

ueroutes = dict(ueRoutingInfo={})

# ueroutes = {
#     "ueRoutingInfo": {
#         "east": {
#             "members": [
#                 'imsi-208930000000007'
#             ],
#             "topology": [
#                 {"A": "gNB1", "B": "UPF-R1"},
#                 {"A": "UPF-R1", "B": "UPF-T1"},
#                 {"A": "UPF-T1", "B": "UPF-C1"}
#             ],
#             "specificPath": [
#                 {"dest": "60.61.0.0/16", "path": ["UPF-R1", "UPF-T1", "UPF-C1"]}
#             ]
#         },
#         "west": {
#             "members": [
#                 'imsi-208930000000001',
#                 'imsi-208930000000003'
#             ],
#             "topology":[
#                 {"A": "gNB1", "B": "UPF-R1"},
#                 {"A": "UPF-R1", "B": "UPF-T2"},
#                 {"A": "UPF-T2", "B": "UPF-C4"}
#             ],
#             "specificPath": [
#                 {"dest": "60.64.0.0/16", "path": ["UPF-R1", "UPF-T2", "UPF-C4"]}
#             ]
#         },
#         "middle": {
#             "members": [
#                 'imsi-208930000000004',
#                 'imsi-208930000000005',
#                 'imsi-208930000000002'
#             ],
#             "topology": [
#                 {"A": "gNB1", "B": "UPF-R1"},
#                 {"A": "UPF-R1", "B": "UPF-T2"},
#                 {"A": "UPF-T2", "B": "UPF-C2"}
#             ],
#             "specificPath": [
#                 {"dest": "60.62.0.0/16", "path": ["UPF-R1", "UPF-T2", "UPF-C2"]}
#             ]
#         }
#     }
# }


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


@proxy.route('/ue-routes/<group_name>', methods=['POST'])
def ueroutes_group_create(group_name):
    try:
        global ueroutes
        if ueroutes['ueRoutingInfo'].get(group_name) is not None:
            response = flask.jsonify({'error': 'group "%s" exists' % group_name})
            response.status_code = 409
            return response

        ueroutes['ueRoutingInfo'][group_name] = {}
        print(ueroutes)
        return ('OK', 200)

    except Exception as e:
        response = flask.jsonify({'error': '%s' % str(e)})
        response.status_code = 500
        print(response)
        return response


@proxy.route('/ue-routes/<group_name>', methods=['GET'])
def ueroutes_group_get(group_name):
    try:
        values = ueroutes['ueRoutingInfo'][group_name]
        response = flask.jsonify(values)
        response.status_code = 200

    except KeyError as e:
        response = flask.jsonify({'error missing key': '%s' % str(e)})
        response.status_code = 404

    except Exception as e:
        response = flask.jsonify({'error': '%s' % str(e)})
        response.status_code = 500

    print(response)
    return response


@proxy.route('/ue-routes/<group_name>/members/<member_name>', methods=['POST'])
def ueroutes_member_create(group_name, member_name):
    try:
        global ueroutes
        members = ueroutes['ueRoutingInfo'][group_name].get('members', [])
        if member_name in members:
            response = flask.jsonify(
                {'error': 'member "%s" exists in group "%s"' % (member_name, group_name)})
            response.status_code = 409
            return response

        ueroutes['ueRoutingInfo'][group_name].setdefault('members', []).append(member_name)
        print(ueroutes)
        return ('OK', 200)
 
    except KeyError as e:
        response = flask.jsonify({'error missing key': '%s' % str(e)})
        response.status_code = 404

    except Exception as e:
        response = flask.jsonify({'error': '%s' % str(e)})
        response.status_code = 500

    print(response)
    return response


@proxy.route('/ue-routes/<group_name>/members', methods=['GET'])
def ueroutes_member_get(group_name):
    try:
        values = ueroutes['ueRoutingInfo'][group_name]['members']
        response = flask.jsonify(values)
        response.status_code = 200
 
    except KeyError as e:
        response = flask.jsonify({'error missing key': '%s' % str(e)})
        response.status_code = 404

    except Exception as e:
        response = flask.jsonify({'error': '%s' % str(e)})
        response.status_code = 500

    print(response)
    return response


@proxy.route('/ue-routes/<group_name>/topology', methods=['POST'])
def ueroutes_topology_create(group_name):
    try:
        values = getMessagePayload()

        global ueroutes
        ueroutes['ueRoutingInfo'][group_name].update(values)
        print(ueroutes)
        return ('OK', 200)
 
    except KeyError as e:
        response = flask.jsonify({'error missing key': '%s' % str(e)})
        response.status_code = 404

    except Exception as e:
        response = flask.jsonify({'error': '%s' % str(e)})
        response.status_code = 500

    print(response)
    return response


@proxy.route('/ue-routes/<group_name>/topology', methods=['GET'])
def ueroutes_topology_get(group_name):
    try:
        values_t = ueroutes['ueRoutingInfo'][group_name]['topology']
        values_s = ueroutes['ueRoutingInfo'][group_name]['specificPath']
        response = flask.jsonify(dict(topology=values_t, specificPath=values_s))
        response.status_code = 200

    except KeyError as e:
        response = flask.jsonify({'error missing key': '%s' % str(e)})
        response.status_code = 404

    except Exception as e:
        response = flask.jsonify({'error': '%s' % str(e)})
        response.status_code = 500

    print(response)
    return response


@proxy.route("/ue-routes", methods=['GET'])
def ueroutes_get():
    '''
    This endpoint is consumed by SMF
    '''
    sys.stdout.write ('Enter /ur-routes\n')
    global prefix
    global ueroutes
    if not ueroutes:
        response = flask.jsonify({'NOT_FOUND': 404})
        response.status_code = 404
        return response
    else:
        ueroutes.update(prefix)
        # important to append prefix for smf not to fail
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
