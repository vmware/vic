import os

BUSYBOX = 'harbor.ci.drone.local/library/busybox' if (os.environ.has_key("DRONE_BUILD_NUMBER") and (int(os.environ['DRONE_BUILD_NUMBER']) != 0)) else 'busybox'
