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
sed -i.bak 's/hostname = reg.mydomain.com/hostname = harbor.longevity/g' harbor.cfg
if [[ ! $(grep harbor.longevity /etc/hosts) ]]; then
    echo "Adding harbor.longevity to /etc/hosts"
    sudo sh -c 'echo "127.0.0.1  harbor.longevity" >> /etc/hosts'
fi

echo "Installing & starting Harbor"
sudo ./install.sh
popd
popd

echo "Preparing Harbor..."
echo "Logging in..."
docker login harbor.longevity --username=admin --password="Harbor12345"
echo "Pulling some images to put in Harbor and putting them in Harbor.."

pushd /home/$USER/vic/tests/resources
for image in $(python -c "vars=__import__('dynamic-vars'); print(\" \".join(vars.images))"); do
    docker pull $image
    docker tag $image harbor.longevity/library/${image}
    docker push harbor.longevity/library/${image}
done
popd
