#!/bin/bash
# Instalar Docker
apt-get update && apt-get install -y docker.io
systemctl start docker
systemctl enable docker
    
# Ejecutar CockroachDB
docker run -d --name=cockroachdb \
  --hostname=cockroachdb-$(hostname) \
  --net=host \
  -v /var/lib/cockroach:/cockroach/data \
  cockroachdb/cockroach:v22.2.8 start \
  --insecure \
  --advertise-addr=$(hostname) \
  --join=COCKROACHDB_JOIN_PLACEHOLDER
