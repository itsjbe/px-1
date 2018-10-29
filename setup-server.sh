#!/bin/bash

locale-gen cs_CZ.UTF-8

export IP_ADDRESS=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1')
export EXPECT_NODES=1

# Setup nomad

apt-get update
apt-get install -y unzip dnsmasq

curl -SLO https://releases.hashicorp.com/nomad/0.8.6/nomad_0.8.6_linux_amd64.zip
unzip nomad_0.8.6_linux_amd64.zip
mv nomad /usr/local/bin/

mkdir -p /var/lib/nomad
mkdir -p /etc/nomad

rm nomad_0.8.6_linux_amd64.zip

cat > server.hcl <<EOF
addresses {
    rpc  = "ADVERTISE_ADDR"
    serf = "ADVERTISE_ADDR"
}
advertise {
    http = "ADVERTISE_ADDR:4646"
    rpc  = "ADVERTISE_ADDR:4647"
    serf = "ADVERTISE_ADDR:4648"
}
bind_addr = "0.0.0.0"
data_dir  = "/var/lib/nomad"
log_level = "DEBUG"
server {
    enabled = true
    bootstrap_expect = EXPECT_NODES
}
EOF

sed -i "s/ADVERTISE_ADDR/${IP_ADDRESS}/" server.hcl
sed -i "s/EXPECT_NODES/${EXPECT_NODES}/" server.hcl
mv server.hcl /etc/nomad/server.hcl

cat > nomad.service <<'EOF'
[Unit]
Description=Nomad
Documentation=https://nomadproject.io/docs/
[Service]
ExecStart=/usr/local/bin/nomad agent -config /etc/nomad
ExecReload=/bin/kill -HUP $MAINPID
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target
EOF

mv nomad.service /etc/systemd/system/nomad.service

systemctl enable nomad
systemctl start nomad


# Setup consul

mkdir -p /var/lib/consul

curl -SLO https://releases.hashicorp.com/consul/1.3.0/consul_1.3.0_linux_amd64.zip
unzip consul_1.3.0_linux_amd64.zip
mv consul /usr/local/bin/consul
rm consul_1.3.0_linux_amd64.zip

cat > consul.service <<'EOF'
[Unit]
Description=consul
Documentation=https://consul.io/docs/
[Service]
ExecStart=/usr/local/bin/consul agent \
  -advertise=ADVERTISE_ADDR \
  -bind=0.0.0.0 \
  -bootstrap-expect=EXPECT_NODES \
  -client=0.0.0.0 \
  -data-dir=/var/lib/consul \
  -server \
  -ui
  
ExecReload=/bin/kill -HUP $MAINPID
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target
EOF

sed -i "s/ADVERTISE_ADDR/${IP_ADDRESS}/" consul.service
sed -i "s/EXPECT_NODES/${EXPECT_NODES}/" consul.service

mv consul.service /etc/systemd/system/consul.service
systemctl enable consul
systemctl start consul



# # Setup Vault

curl -SLO https://releases.hashicorp.com/vault/0.11.4/vault_0.11.4_linux_amd64.zip
unzip vault_0.11.4_linux_amd64.zip
mv vault /usr/local/bin/vault
rm vault_0.11.4_linux_amd64.zip

mkdir -p /etc/vault

cat > /etc/vault/vault.hcl <<'EOF'
disable_mlock = true
backend "consul" {
  advertise_addr = "http://ADVERTISE_ADDR:8200"
  address = "127.0.0.1:8500"
  path = "vault"
}
listener "tcp" {
  address = "ADVERTISE_ADDR:8200"
  tls_disable = 1
}
EOF

sed -i "s/ADVERTISE_ADDR/${IP_ADDRESS}/" /etc/vault/vault.hcl

cat > /etc/systemd/system/vault.service <<'EOF'
[Unit]
Description=Vault
Documentation=https://vaultproject.io/docs/
[Service]
ExecStart=/usr/local/bin/vault server \
  -config /etc/vault/vault.hcl
  
ExecReload=/bin/kill -HUP $MAINPID
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target
EOF

# systemctl enable vault
# systemctl start vault


# Setup dnsmasq

mkdir -p /etc/dnsmasq.d
cat > /etc/dnsmasq.d/10-consul <<'EOF'
server=/consul/127.0.0.1#8600
EOF

systemctl enable dnsmasq
systemctl start dnsmasq
