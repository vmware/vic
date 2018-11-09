# Building:
# docker build --no-cache -t vic-local-repo-centos-6.9 -f infra/integration-image/local-repo/Dockerfile.centos-6.9 infra/integration-image/local-repo/
# docker tag vic-local-repo-centos-6.9 gcr.io/eminent-nation-87317/vic-local-repo:centos-6.9
# gcloud auth login
# gcloud docker -- push gcr.io/eminent-nation-87317/vic-local-repo:centos-6.9
#
# Running:
# docker run -d -p 80:80 vic-local-repo-centos-6.9
FROM fedora:21

RUN yum install yum-plugin-ovl -y && yum install wget createrepo nginx -y

RUN mkdir -p /usr/share/nginx/html/centos/Packages && \
    mkdir -p /usr/share/nginx/html/centos-updates/Packages

ENV EXCLUDE_LIST "index.html*,openjdk*,openjre*,postgresql*,python-*,\
python3*,ruby*,subversion*,gnome*,NetworkManager*,cloud*,docker*,ktap*,\
kubernetes*,linux-docs*,linux-sound*,linux-tools*,docbook*,httpd*,go-*,jna*,\
linux-debuginfo*,linux-dev*,linux-docs*,linux-drivers*,linux-oprofile*,linux-sound*,\
linux-tools*,linux-esx-debuginfo*,linux-esx-devel*,linux-esx-docs*,nginx*,sysdig*"

RUN wget http://mirror.centos.org/centos/6/os/x86_64/RPM-GPG-KEY-CentOS-6 -O /usr/share/nginx/html/centos/RPM-GPG-KEY-CentOS-6
RUN wget -e robots=off -r -nH -nd -np -R $EXCLUDE_LIST http://mirror.centos.org/centos/6/os/x86_64/Packages/ -P /usr/share/nginx/html/centos/Packages/
RUN wget -e robots=off -r -nH -nd -np -R $EXCLUDE_LIST http://mirror.centos.org/centos/6/updates/x86_64/Packages/ -P /usr/share/nginx/html/centos-updates/Packages/

RUN createrepo /usr/share/nginx/html/centos
RUN createrepo /usr/share/nginx/html/centos-updates

RUN echo "daemon off;" >> /etc/nginx/nginx.conf

EXPOSE 80

CMD [ "/usr/sbin/nginx" ]
