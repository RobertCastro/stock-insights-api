# Crear cuenta de servicio para GKE
resource "google_service_account" "gke_sa" {
  account_id   = "gke-service-account"
  display_name = "GKE Service Account"
}

resource "google_project_iam_member" "gke_sa_roles" {
  for_each = toset([
    "roles/container.developer",
    "roles/storage.objectViewer"
  ])
  
  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.gke_sa.email}"
}

# Módulo de red
module "network" {
  source = "./modules/network"
  
  network_name        = "backend-network"
  region              = var.region
  private_subnet_cidr = "10.0.0.0/16"
}

# Módulo de CockroachDB
module "cockroachdb" {
  source = "./modules/cockroachdb"
  
  instance_count        = var.cockroachdb_instance_count
  machine_type          = var.cockroachdb_machine_type
  zone                  = var.zone
  disk_size_gb          = var.cockroachdb_disk_size_gb
  subnetwork_id         = module.network.private_subnet_id
  service_account_email = google_service_account.gke_sa.email
}

# Módulo de Kubernetes
module "kubernetes" {
  source = "./modules/kubernetes"
  
  cluster_name    = "backend-cluster"
  zone            = var.zone
  network_id      = module.network.vpc_id
  subnet_id       = module.network.private_subnet_id
  node_count      = var.gke_node_count
  machine_type    = var.gke_machine_type
  service_account = google_service_account.gke_sa.email
  
  depends_on = [module.network, module.cockroachdb]
}