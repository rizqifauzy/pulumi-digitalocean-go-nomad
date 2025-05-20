pulumi new go

pulumi config set digitalocean:token {DIGITAL OCEAN TOKEN} --secret


ssh-keygen -f dostart-keypair

nomad tls cert create -server -region global

pulumi config set publicKeyPath dostart-keypair.pub
pulumi config set privateKeyPath dostart-keypair.key

pulumi stack output dropletIP


nomad node status \
    -ca-cert=nomad-agent-ca.pem \
    -client-cert=global-cli-nomad.pem \
    -client-key=global-cli-nomad-key.pem \
    -address=https://{IP_ADDRESS}:4646

export NOMAD_ADDR=https://127.0.0.1:4646

nomad status

nomad acl bootstrap

export NOMAD_TOKEN="8640c485-ab97-08c3-3b60-b45cae48740f"
