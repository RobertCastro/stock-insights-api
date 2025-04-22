resource "google_container_cluster" "primary" {
  name     = var.cluster_name
  location = var.zone
  
  # Cluster de zona única
  initial_node_count = 1
  
  # Red y subred específicas
  network    = var.network_id
  subnetwork = var.subnet_id
  
  # Desactivar el pool de nodos por defecto
  remove_default_node_pool = true
  
  # Habilitar Private Cluster
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = "172.16.0.0/28"
  }
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "${var.cluster_name}-node-pool"
  location   = var.zone
  cluster    = google_container_cluster.primary.name
  node_count = var.node_count
  
  node_config {
    machine_type = var.machine_type
    
    service_account = var.service_account
    oauth_scopes    = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}