# Información de red
output "vpc_id" {
  description = "ID de la VPC creada"
  value       = module.network.vpc_id
}

output "subnet_id" {
  description = "ID de la subred privada"
  value       = module.network.private_subnet_id
}

# Información de CockroachDB
output "cockroachdb_instances" {
  description = "Nombres de las instancias de CockroachDB"
  value       = module.cockroachdb.instance_names
}

output "cockroachdb_internal_ips" {
  description = "IPs internas de las instancias de CockroachDB"
  value       = module.cockroachdb.internal_ips
}

# Información de Kubernetes
output "kubernetes_cluster_name" {
  description = "Nombre del cluster de Kubernetes"
  value       = module.kubernetes.cluster_name
}

output "kubernetes_endpoint" {
  description = "Endpoint del cluster de Kubernetes"
  value       = module.kubernetes.endpoint
  sensitive   = true
}

output "kubernetes_service_account" {
  description = "Cuenta de servicio utilizada por el cluster"
  value       = google_service_account.gke_sa.email
}

# Información para acceso al backend
output "backend_service_ip" {
  description = "IP del servicio del backend"
  value       = "Ejecuta 'kubectl get service backend-service -o jsonpath='{.status.loadBalancer.ingress[0].ip}'' para obtener la IP"
}