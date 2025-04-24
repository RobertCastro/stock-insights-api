output "kubernetes_cluster_name" {
  value       = google_container_cluster.primary.name
  description = "Nombre del cluster de GKE"
}

output "kubernetes_cluster_host" {
  value       = google_container_cluster.primary.endpoint
  description = "Host del cluster de GKE"
  sensitive   = true
}

output "project_id" {
  value       = var.project_id
  description = "ID del proyecto GCP"
}

output "region" {
  value       = var.region
  description = "Región de GCP donde se han desplegado los recursos"
}
