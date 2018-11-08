path "secret/toriiapp" {
  capabilities = ["read", "list"]
}

path "mysql/creds/toriiapp" {
  capabilities = ["read", "list"]
}

path "sys/renew/*" {
  capabilities = ["update"]
}


