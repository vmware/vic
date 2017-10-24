# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

#!/bin/bash
set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 harbor-version"
    exit 1
fi
version=$1
pushd /home/$USER
[ -e harbor ] \
    && echo "/home/$USER/harbor exists. Delete it if you want to install a newer version and re-run $0" \
    && pushd harbor && docker-compose start && popd && exit 0

echo "Pulling down version ${version} of Harbor..."
wget https://github.com/vmware/harbor/releases/download/v${version}/harbor-online-installer-v${version}.tgz -qO - | tar xz
pushd harbor
echo "Configuring Harbor"
sed -i 's/hostname = reg.mydomain.com/hostname = willie.eng.vmware.com/g' harbor.cfg


if [ ! -e /home/$USER/harbor-keys ]; then
    echo "Generating SSL cert for Harbor"
    mkdir -p ~/harbor-keys
    pushd ~/harbor-keys
    openssl genrsa -des3 -passout pass:x -out harbor.pass.key 2048
    openssl rsa -passin pass:x -in harbor.pass.key -out harbor.key
    rm harbor.pass.key
    openssl req -new -key harbor.key -out harbor.csr \
            -subj "/C=US/ST=TX/L=Austin/O=CNAU/OU=Testing/CN=willie.eng.vmware.com"
    openssl x509 -req -days 365 -in harbor.csr -signkey harbor.key -out harbor.crt
    popd
fi


sed -i "s#/data/cert/server.crt#/home/$USER/harbor-keys/harbor.crt#g" harbor.cfg
sed -i "s#/data/cert/server.key#/home/$USER/harbor-keys/harbor.key#g" harbor.cfg

echo "Installing & starting Harbor"
sudo ./install.sh
popd
popd

echo "Preparing Harbor..."
echo "Logging in..."
docker login willie.eng.vmware.com --username=admin --password="Harbor12345"
echo "Pulling some images to put in Harbor and putting them in Harbor.."

pushd /home/$USER/vic/tests/resources
for image in $(python -c "vars=__import__('dynamic-vars'); print(\" \".join(vars.images))"); do
    docker pull $image
    docker tag $image willie.eng.vmware.com/library/${image}
    docker push willie.eng.vmware.com/library/${image}
done
popd
