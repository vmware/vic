echo "Removing VIC directory if present"
rm -rf vic

input=$(wget -O - https://vmware.bintray.com/vic-repo |tail -n5 |head -n1 |cut -d':' -f 2 |cut -d'.' -f 3| cut -d'>' -f 2)

echo "Downloading bintray file"
wget https://vmware.bintray.com/vic-repo/$input.tar.gz

echo "Extracting .tar.gz"
tar xzf $input.tar.gz

echo "Deleting .tar.gz vic file"
rm $input.tar.gz

drone exec --trusted -e test="pybot tests/manual-test-cases/Group5-Functional-Tests/5-1-Distributed-Switch.robot" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

mv output.xml output_51.xml
mv log.html log_51.html

drone exec --trusted -e test="pybot tests/manual-test-cases/Group5-Functional-Tests/5-2-Cluster.robot" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

mv output.xml output_52.xml
mv log.html log_52.html

drone exec --trusted -e test="pybot tests/manual-test-cases/Group5-Functional-Tests/5-4-High-Availability.robot" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

mv output.xml output_54.xml
mv log.html log_54.html

drone exec --trusted -e test="pybot tests/manual-test-cases/Group5-Functional-Tests/5-5-Heterogenous-ESXi.robot" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

mv output.xml output_55.xml
mv log.html log_55.html

drone exec --trusted -e test="pybot tests/manual-test-cases/Group5-Functional-Tests/5-6-VSAN.robot" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

mv output.xml output_56.xml
mv log.html log_56.html

drone exec --trusted -e test="pybot tests/manual-test-cases/Group5-Functional-Tests/5-7-NSX.robot" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

mv output.xml output_57.xml
mv log.html log_57.html

drone exec --trusted -e test="pybot tests/manual-test-cases/Group5-Functional-Tests/5-8-DRS.robot" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

mv output.xml output_58.xml
mv log.html log_58.html

drone exec --trusted -e test="pybot tests/manual-test-cases/Group5-Functional-Tests/5-9-Private-Registry.robot" -E nightly_test_secrets.yml --yaml .drone.nightly.yml

drone exec --trusted -e test="sh tests/Nightly/Upload-logs.sh $input" -E nightly_test_secrets.yml --yaml .drone.nightly.yml
