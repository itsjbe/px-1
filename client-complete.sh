#!/bin/bash


echo 'datasource_list: [ None ]' | sudo -s tee /etc/cloud/cloud.cfg.d/90_dpkg.cfg

apt-get purge cloud-init -y

rm -rf /etc/cloud/; rm -rf /var/lib/cloud/

systemctl show -p WantedBy network-online.target

systemctl disable open-iscsi.service
systemctl disable iscsid.service

apt remove open-iscsi -y

apt install unzip curl iputils-ping vim -y

apt autoremove -y