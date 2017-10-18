import os

IS_LOCAL = False if (os.environ.has_key("DRONE_BUILD_NUMBER") and (int(os.environ['DRONE_BUILD_NUMBER']) != 0)) else True

images = ['busybox', 'alpine', 'nginx','debian', 'ubuntu', 'redis']

for image in images:
    name = image if IS_LOCAL else 'harbor.ci.drone.local/library/{}'.format(image)
    exec("{} = '{}'".format(image.upper(), name))
