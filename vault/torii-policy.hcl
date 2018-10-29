path "secret/torii" {
  capabilities = ["read", "list"]
}

path "mysql/creds/torii" {
  capabilities = ["read", "list"]
}

path "sys/renew/*" {
  capabilities = ["update"]
}


