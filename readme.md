## Bootstrap servers 1 and clients 1, 2
```
./setup-ns.sh
```

## Create database

```
ssh jbe@mysql
mysql -u root -p
```
```
mysql> CREATE DATABASE toriiapp;
```

## Complete the setup of the vault

set hosts record for ns-1 to correct eth address

``` bash
export VAULT_ADDR=http://ns-1:8200
vault operator init
vault operator unseal
vault status
vault login <root-token>
```

```bash
vault secrets enable mysql
```
```bash
vault write mysql/config/connection connection_url="root:q@tcp(mysql:3306)/"
```
```bash
vault write mysql/roles/toriiapp sql="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT ALL PRIVILEGES ON toriiapp.* TO '{{name}}'@'%';"
```






### Create the Toriiapp Policy and Token

```bash
vault policy write torii vault/torii-policy.hcl
```

```bash

vault token create \
  -policy="toriiapp" \
  -display-name="toriiapp"
```

### Create the Toriiapp Secret

```
vault read mysql/creds/torii
```

```
vault write secret/torii jwtsecret=using-some-kind-of-secret-to-sign-jwt
```

