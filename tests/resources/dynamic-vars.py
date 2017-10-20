import os
from enum import Enum


class TestEnvironment(Enum):
    LOCAL = 0
    DRONE = 1
    LONGEVITY = 2


def getEnvironment():
    if (os.environ.has_key("DRONE_BUILD_NUMBER") and (int(os.environ['DRONE_BUILD_NUMBER']) != 0)):
        return TestEnvironment.DRONE
    elif os.environ.has_key("LONGEVITY"):
        return TestEnvironment.LONGEVITY
    else:
        return TestEnvironment.LOCAL

def getName(image):
    return {TestEnvironment.DRONE: 'harbor.ci.drone.local/library/{}'.format(image),
            TestEnvironment.LONGEVITY: 'harbor.longevity/library/{}'.format(image),
            TestEnvironment.LOCAL: image}[getEnvironment()]

# this global variable (images) is used by the Longevity scripts. If you change this, change those!
# and don't inline it!
images = ['busybox', 'alpine', 'nginx','debian', 'ubuntu', 'redis']
for image in images:
    exec("{} = '{}'".format(image.upper().replace(':', '_').replace('.', '_'), getName(image)))
