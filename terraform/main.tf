provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_project_service" "services" {
  for_each = toset([
    "container.googleapis.com",     # GKE
    "artifactregistry.googleapis.com", # Artifact Registry
    "dns.googleapis.com",           # Cloud DNS
    "compute.googleapis.com",       # Compute Engine
    "certificatemanager.googleapis.com" # Certificate Manager
  ])
  project = var.project_id
  service = each.key

  disable_dependent_services = true
}

# Crear el repositorio de Artifact Registry si no existe
resource "google_artifact_registry_repository" "docker_repo" {
  provider = google
  location = var.region
  repository_id = "docker-images"
  description = "Docker repository for Stock Insights API"
  format = "DOCKER"
  depends_on = [google_project_service.services]
}

# Red VPC para el cluster
resource "google_compute_network" "vpc_network" {
  name                    = "stock-insights-vpc"
  auto_create_subnetworks = false
  depends_on              = [google_project_service.services]
}

# Subred para el cluster
resource "google_compute_subnetwork" "subnet" {
  name          = "stock-insights-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.region
  network       = google_compute_network.vpc_network.id
  
  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "10.1.0.0/16"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "10.2.0.0/16"
  }
}

# Firewall para permitir acceso al cluster
resource "google_compute_firewall" "allow_internal" {
  name    = "allow-internal"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  source_ranges = ["10.0.0.0/16", "10.1.0.0/16", "10.2.0.0/16"]
}

# Cluster de GKE
resource "google_container_cluster" "primary" {
  name     = var.cluster_name
  location = var.region
  
  # Eliminar el default node pool
  remove_default_node_pool = true
  initial_node_count       = 1

  network    = google_compute_network.vpc_network.id
  subnetwork = google_compute_subnetwork.subnet.id

  ip_allocation_policy {
    cluster_secondary_range_name  = "pod-ranges"
    services_secondary_range_name = "services-range"
  }

  # Habilitar la autenticación de workload identity
  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  # Habilitar Network Policy para mejor seguridad
  network_policy {
    enabled = true
  }

  # Configuración de certificado SSL
  addons_config {
    http_load_balancing {
      disabled = false
    }
  }

  depends_on = [google_project_service.services]
}

# Node pool para las aplicaciones
resource "google_container_node_pool" "primary_nodes" {
  name       = "${var.cluster_name}-node-pool"
  location   = var.region
  cluster    = google_container_cluster.primary.name
  node_count = 2

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/compute",
    ]

    labels = {
      env = "production"
    }

    machine_type = "e2-medium"
    disk_size_gb = 50
    disk_type    = "pd-standard"
    
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }
}

# Certificado SSL administrado por Google
resource "google_certificate_manager_certificate" "default" {
  name        = "api-soyrobert-co-cert"
  description = "Certificate for api.soyrobert.co"
  
  managed {
    domains = ["api.soyrobert.co"]
  }
  
  depends_on = [google_project_service.services]
}