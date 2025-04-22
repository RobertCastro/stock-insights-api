output "instance_names" {
  description = "Nombres de las instancias de CockroachDB"
  value       = google_compute_instance.cockroachdb[*].name
}

output "internal_ips" {
  description = "IPs internas de las instancias de CockroachDB"
  value       = google_compute_instance.cockroachdb[*].network_interface[0].network_ip
}