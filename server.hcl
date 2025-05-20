name = "nomad"

data_dir = "/var/lib/nomad"

bind_addr = "0.0.0.0"

advertise {
  http = "167.172.71.77"
  rpc  = "167.172.71.77"
  serf = "167.172.71.77"
}

ports {
  http = 4646
  rpc  = 4647
  serf = 4648
}

tls {
  http = true
  rpc  = true

  ca_file   = "/etc/certs/ca.crt"
  cert_file = "/etc/certs/nomad.crt"
  key_file  = "/etc/certs/nomad.key"
}

datacenter = "dc1"

log_level = "INFO"
log_file  = "/etc/nomad.d/nomad.log"

server {
  enabled          = true
  bootstrap_expect = 1
  encrypt          = "YHYyBkYOUtMIvEOdvJ1Cx2AeWPbv0y4HA4UwEB35avY="
}

client {
  enabled = true
}

acl {
  enabled    = true
  token_ttl  = "30s"
  policy_ttl = "60s"
  role_ttl   = "60s"
}