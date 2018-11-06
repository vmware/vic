#!/usr/bin/python

import json
import requests


class NSXtClient(object):
    _VERB_RESP_CODES = {
        'get': [requests.codes.ok],
        'post': [requests.codes.created, requests.codes.ok],
        'put': [requests.codes.ok],
        'delete': [requests.codes.ok]
    }
    _DEFAULT_HEADERS = {'Accept': 'application/json',
                        'Content-Type': 'application/json'}

    def __init__(self, nsxt_uri, nsxt_user, nsxt_password, insecure=True,
                 ca_file=None, default_headers=_DEFAULT_HEADERS):
        session = requests.Session()
        session.auth = (nsxt_user, nsxt_password)
        session.max_redirects = 0

        session.verify = not insecure
        if session.verify and ca_file:
            session.verify = ca_file

        self._conn = session
        if nsxt_uri.startswith('https://'):
            self._url_prefix = "%s/api/v1/" % nsxt_uri
        else:
            self._url_prefix = "https://%s/api/v1/" % nsxt_uri
        self._default_headers = default_headers

    def url_get(self, url, params=None):
        return self._rest_call(url, method='GET', params=params)

    def url_delete(self, url):
        return self._rest_call(url, method='DELETE')

    def url_put(self, url, body):
        return self._rest_call(url, method='PUT', body=body)

    def url_post(self, url, body):
        return self._rest_call(url, method='POST', body=body)

    def _validate_result(self, result, expected, operation):
        if result.status_code not in expected:
            result_msg = result.json() if result.content else ''
            if type(result_msg) is dict:
                result_msg = result_msg.get('error_message', result_msg)
            raise Exception(
                'Requested operation %s \nUnexpected response %s received'
                'with msg %s' % (operation, result.status_code, result_msg))

    def _build_url(self, uri):
        return self._url_prefix + uri

    def _rest_call(self, url, method='GET', body=None, params=None):
        if body is not None:
            body = json.dumps(body)
        request_headers = {}
        request_headers.update(self._default_headers)
        request_url = self._build_url(url)

        do_request = getattr(self._conn, method.lower())

        result = do_request(request_url, data=body,
                            headers=request_headers, params=params)

        self._validate_result(result,
                              NSXtClient._VERB_RESP_CODES[method.lower()],
                              ("%(verb)s %(url)s") % {'verb': method, 'url': request_url})
        return result.json() if result.content else result


def get_overlay_transport_zone(nsxt_uri, nsxt_username, nsxt_password):
    nsxclient = NSXtClient(nsxt_uri, nsxt_username, nsxt_password)
    result = nsxclient.url_get('transport-zones')
    for trans in result['results']:
        if trans['transport_type'] == "OVERLAY":
            return trans['id']


def delete_nsxt_logical_switch(nsxt_uri, nsxt_username, nsxt_password, logical_switch_id):
    nsxclient = NSXtClient(nsxt_uri, nsxt_username, nsxt_password)
    uri = "logical-switches/%s" % logical_switch_id
    nsxclient.url_delete(uri)


def get_nsxt_overlay_logical_switch(nsxt_uri, nsxt_username, nsxt_password,
                                    transport_zone_id, logical_switch_name):
    nsxclient = NSXtClient(nsxt_uri, nsxt_username, nsxt_password)
    result = nsxclient.url_get('logical-switches')
    for ls in result['results']:
        if (ls['display_name'] == logical_switch_name and
                    ls['transport_zone_id'] == transport_zone_id):
            return ls


def create_nsxt_logical_switch(nsxt_uri, nsxt_username, nsxt_password, logical_switch_name):
    # As there is no transport zone name returned, so we used first returned overlay TZ
    transport_zone_id = get_overlay_transport_zone(nsxt_uri, nsxt_username, nsxt_password)
    if transport_zone_id is None:
        return
    existing_ls = get_nsxt_overlay_logical_switch(nsxt_uri, nsxt_username, nsxt_password,
                                                  transport_zone_id, logical_switch_name)
    if existing_ls:
        return existing_ls['id']

    nsxclient = NSXtClient(nsxt_uri, nsxt_username, nsxt_password)
    body = {
        'transport_zone_id': transport_zone_id,
        'replication_mode': 'MTEP',
        'admin_state': 'UP',
        'display_name': logical_switch_name
    }
    result = nsxclient.url_post('logical-switches', body)
    return result['id']
