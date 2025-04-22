resource "local_file" "startup_script" {
  content  = <<-EOT
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
  EOT
  filename = "${path.module}/startup-script.tpl"
}

# Luego utilizamos las instancias
resource "google_compute_instance" "cockroachdb" {
  count        = var.instance_count
  name         = "cockroachdb-${count.index}"
  machine_type = var.machine_type
  zone         = var.zone
  
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
      size  = var.disk_size_gb
      type  = "pd-ssd"
    }
  }
  
  network_interface {
    subnetwork = var.subnetwork_id
    # Sin IP externa - acceso solo a través de la red privada
  }
  
  # Modificamos el script según la instancia
  metadata_startup_script = replace(
    local_file.startup_script.content,
    "COCKROACHDB_JOIN_PLACEHOLDER",
    count.index > 0 ? "cockroachdb-0:26257" : ""
  )
  
  service_account {
    email  = var.service_account_email
    scopes = ["cloud-platform"]
  }
  
  tags = ["cockroachdb", "database"]
  
  depends_on = [local_file.startup_script]
}