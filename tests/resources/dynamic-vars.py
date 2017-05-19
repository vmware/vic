import os

IS_LOCAL = False if (os.environ.has_key("DRONE_BUILD_NUMBER") and (int(os.environ['DRONE_BUILD_NUMBER']) != 0)) else True

BUSYBOX =  'busybox' if IS_LOCAL else 'harbor.ci.drone.local/library/busybox'
ALPINE =  'alpine' if IS_LOCAL else 'harbor.ci.drone.local/library/alpine'
NGINX =  'nginx' if IS_LOCAL else 'harbor.ci.drone.local/library/nginx'
DEBIAN =  'debian' if IS_LOCAL else 'harbor.ci.drone.local/library/debian'
UBUNTU =  'ubuntu' if IS_LOCAL else 'harbor.ci.drone.local/library/ubuntu'
REDIS =  'redis' if IS_LOCAL else 'harbor.ci.drone.local/library/redis'
