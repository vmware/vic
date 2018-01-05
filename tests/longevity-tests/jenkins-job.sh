#!/bin/bash
cd ~/vic

git clean -fd
git fetch https://github.com/vmware/vic master
git pull

cp ~/secrets .
tests/longevity-tests/run-longevity.bash $1
id=`docker ps -lq`
echo $id

docker logs -f $id

docker cp $id:/tmp $id
tar -cvzf $id.tar.gz $id
gsutil cp $id.tar.gz gs://vic-longevity-results/

echo $id
rc=`docker inspect --format='{{.State.ExitCode}}' $id`
exit $rc
