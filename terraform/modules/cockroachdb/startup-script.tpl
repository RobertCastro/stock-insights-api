#!/bin/bash
# Instalar Docker
apt-get update && apt-get install -y docker.io
systemctl start docker
systemctl enable docker
    
# Deshabilitar firewall para pruebas
systemctl stop ufw || true
systemctl disable ufw || true
iptables -F || true
    
# Verificar IP interna
INTERNAL_IP=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip)
echo "Internal IP: $INTERNAL_IP" > /tmp/ip_debug.log
    
# Ejecutar CockroachDB con configuración más permisiva
docker run -d --name=cockroachdb \
  --hostname=cockroachdb-$(hostname) \
  --net=host \
  -v /var/lib/cockroach:/cockroach/data \
  -p 26257:26257 \
  cockroachdb/cockroach:v22.2.8 start \
  --insecure \
  --listen-addr=0.0.0.0:26257 \
  --advertise-addr=$INTERNAL_IP:26257 \
  --accept-sql-without-tls \
  --http-addr=0.0.0.0:8080 \
  --join=COCKROACHDB_JOIN_PLACEHOLDER
    
# Debug: Verificar que el puerto está abierto
apt-get install -y net-tools lsof netcat
echo "Testing ports:" >> /tmp/ip_debug.log
netstat -tuln >> /tmp/ip_debug.log
docker logs cockroachdb >> /tmp/docker_logs.log
