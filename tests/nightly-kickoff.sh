NIMBUS_USER=
NIMBUS_PASSWORD=
NIMBUS_GW=nimbus-gateway.eng.vmware.com
SLACK_URL=https://hooks.slack.com/services/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
GS_PROJECT_ID=eminent-nation-87317
GS_CLIENT_EMAIL=vic-ci-logs@eminent-nation-87317.iam.gserviceaccount.com
GS_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDmv6+b61zT/Wc/\n2ahOfIrjgKQxhucGizflqp0C520urDGWYB4qd1/W3MGP9z7rjmFVOeT2l6D+OcSw\nn7vvEzjDFce5OOlT1OAQTn1m5fqj7yurEEusMau6LQJpbf0Yv1IbpTI7BMIVjc0l\nvNvWV10mUZUlSHoYEkrMBkJE1vNZO961wjYH/nyWd2VI+Xx5bc4IK28SghPaFffQ\nAN9uTHGXdGaK/rsakfwA/hFEfEpCoaCqIzJTg4gdn6AfsjCXkKCNeMUVW2/xjpKW\nhEI3emgBuFZt3/GsdXGmF2MagVe/35r3oxCNHfDK3sgZUMDr/40RqNa0dgjF6OyD\nfJgeBmkXAgMBAAECggEBALyxE5GVLhvMpJn6Cz/jaFAR6jL19gXL4rHUiwDM7uOz\nu/kUMJbZd23kqARqUvGdRMrExQ9Bf01lQAqPFMe0GD6vmNtGRsde1LuA89spRYS5\nGCSS9s6g76UXGVnNr6KFEUe6FxFcGro1cwThI4RrfKjRHf2W/wCgNLoShC52+BiG\n5v7VpkZJUnrrUtESqd2GaOmOrMyJcyecTVuL1ODlGFtk9Cf4h+CpmTHWAJTBMvPv\nEBP0MjqyCcl8oBktNVzC0z53DvsUNhD6jgvfPZ+AboCebs0xtrtMDaooCvNQ4+2K\nuJoAG9saYyrLeKMULSSKH94+NUUvX2JkDGH2NJoXfgECgYEA/d7NWpGStMkkkI37\noeuZhqyn2n/KLC+aTG4PiWCGdwziul6Q0JlrrwGk+zFz6kNKmGim+mw3JCmHtcYF\n5GJra7CEAV6/e2tnhP2NY8d4wip/N/58xZhWuCSxHur/O8/TTN4tjJ4dIugcf5a3\nk0ZmCtAet2RWHroIqaMVEwAHhPcCgYEA6K86sHPJP/eF5YYdVB0rzDrhcYPSseSS\na/JtOd31F2QZhx654zMlfUjM5z+LVoisf9CntyAabCceGUALOtlGwkoWD8naqPLo\nk6l8bAdiTV5nK5oSGaKb8xC3JHOM2d89dT82NDPWDOCCwQvN4BFUjGtbBmtF6n8d\ncuI578+i1OECgYAidN4MX9u4m+BRmmO/21lQFRkHJ/cJvkBEBWAodihp+h6/ytv+\n5APgkemRimnALvft7a5UKOHnD5fyzPi5wb3wtNmF0hVNLAu12jAZjdZPDDOOJwVK\nUF3cymYb2ytfM9rrAPDPuBoeRcCwdIVgANsStqKko4Ko0vkgBRl0JbnfiwKBgQDD\nix9jUqr6WuXnsgHLwoggJgt3/jR+03xJw34Pd3yVn8XkS+okCcOjuh6Y6EoM+ucc\nsxl/SDdsVKNyzOOjHR3eAazwr85W1WynS3QIxVvTcVZ6ygwUBxfP+Wgv9fuUzYs0\nkV7YGAf24maAHY9ykp3fNAlXJ6emHhV9iqjt5C0PgQKBgASEqRg+lPo0a6r9e/cW\nRwYs/GNvBiZM3V2ttmwz7nyzRuOakcocFk/3ujz58cY2g4UOL/8Rm3ULqf2v+3U1\nDmT12CenkVHB6DE1fngiSUknQ6irf24XdvOiqUVCJAj3eFpmaFcN3NFL/iccGLjq\n8sTizN0defLu1MTOdgDrMzls\n-----END PRIVATE KEY-----\n"

cd "$(git rev-parse --show-toplevel)"

tests=${*#${PWD}/}

input=$(wget -O - https://vmware.bintray.com/vic-repo |tail -n5 |head -n1 |cut -d':' -f 2 |cut -d'.' -f 3| cut -d'>' -f 2)
echo "https://console.cloud.google.com/m/cloudstorage/b/vic-ci-logs/o/functional_logs_"$input".zip?authuser=1"
echo "Downloading bintray file"
wget https://vmware.bintray.com/vic-repo/$input.tar.gz
echo "Removing VIC directory if present"
rm -rf vic
echo "Extracting .tar.gz"
tar xf $input.tar.gz
echo "Deleting .tar.gz vic file"
rm $input.tar.gz

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
    commands:
       - pybot tests/manual-test-cases/Group5-Functional-Tests/5-1-Distributed-Switch.robot
CONFIG
)

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY $GITHUB_AUTOMATION_API_KEY
      DRONE_TOKEN:  $DRONE_TOKEN
    commands:
       - pybot tests/manual-test-cases/Group5-Functional-Tests/5-2-Cluster.robot
CONFIG
)

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
    commands:
       - pybot tests/manual-test-cases/Group5-Functional-Tests/5-4-High-Availability.robot
CONFIG
)

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
    commands:
       - pybot tests/manual-test-cases/Group5-Functional-Tests/5-5-Heterogenous-ESXi.robot
CONFIG
)

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
    commands:
       - pybot tests/manual-test-cases/Group5-Functional-Tests/5-6-VSAN.robot
CONFIG
)

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
    commands:
       - pybot tests/manual-test-cases/Group5-Functional-Tests/5-7-NSX.robot
CONFIG
)

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
    commands:
       - pybot tests/manual-test-cases/Group5-Functional-Tests/5-8-DRS.robot
CONFIG
)

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
    commands:
       - pybot tests/manual-test-cases/Group5-Functional-Tests/5-9-Private-Registry.robot
CONFIG
)

drone exec --trusted --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

build:
  integration-test:
    image: vmware-docker-ci-repo.bintray.io/integration/vic-test:1.8
    pull: true
    environment:
      BIN: vic
      GOPATH: /drone
      SHELL: /bin/bash
      DOCKER_API_VERSION: "1.21"  
      LOG_TEMP_DIR: install-logs     
      TEST_TIMEOUT: $TEST_TIMEOUT
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
      NIMBUS_USER: $NIMBUS_USER
      NIMBUS_PASSWORD: $NIMBUS_PASSWORD
      NIMBUS_GW: $NIMBUS_GW
      SLACK_URL: $SLACK_URL
      GS_PROJECT_ID: $GS_PROJECT_ID
      GS_CLIENT_EMAIL: $GS_CLIENT_EMAIL
      GS_PRIVATE_KEY: $GS_PRIVATE_KEY
    commands:
       - tests/nightly-test.sh 
CONFIG
)

drone exec --trusted --notify --yaml <(cat <<CONFIG
---
clone:
  path: github.com/vmware/vic
  tags: true

notify:
  slack:
    webhook_url: $SLACK_URL
    channel: mwilliamson-staff
    username: drone
    template: >
      build https://ci.vmware.run/vmware/vic/$input finished with a {{ build.status }} status, find the logs here: https://console.cloud.google.com/m/cloudstorage/b/vic-ci-logs/o/functional_logs_$input.zip?authuser=1

CONFIG
)



