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


curl -SLO https://releases.hashicorp.com/nomad/0.8.6/nomad_0.8.6_linux_amd64.zip
curl -SLO https://releases.hashicorp.com/consul/1.3.0/consul_1.3.0_linux_amd64.zip
curl -SLO https://releases.hashicorp.com/vault/0.11.4/vault_0.11.4_linux_amd64.zip
