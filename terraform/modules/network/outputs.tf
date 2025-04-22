output "vpc_id" {
  description = "ID de la VPC creada"
  value       = google_compute_network.vpc.id
}

output "private_subnet_id" {
  description = "ID de la subred privada"
  value       = google_compute_subnetwork.private.id
}