#!/bin/bash -e

apt-get update

# set GOPATH based on shared folder of vagrant
pro="/home/"${BASH_ARGV[0]}"/.profile"
echo "export GOPATH="${BASH_ARGV[1]} >> $pro

# add GOPATH/bin to the PATH
echo "export PATH=$PATH:"${BASH_ARGV[1]}"/bin" >> $pro

packages=(curl lsof strace git shellcheck)

for package in "${packages[@]}" ; do
  apt-get -y install "$package"
done

if [ ! -d "/usr/local/go" ] ; then
  (cd /usr/local &&
   (curl --silent -L https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz | tar -zxf -) &&
   ln -s /usr/local/go/bin/* /usr/local/bin/)
fi

cat << EOF > /etc/systemd/system/docker.service
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network.target docker.socket
Requires=docker.socket

[Service]
Type=notify
ExecStart=/usr/bin/docker daemon -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock -D
MountFlags=slave
LimitNOFILE=1048576
LimitNPROC=1048576
LimitCORE=infinity

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl restart docker
