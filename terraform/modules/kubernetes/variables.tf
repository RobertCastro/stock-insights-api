variable "cluster_name" {
  description = "Nombre del cluster de Kubernetes"
  type        = string
}

variable "zone" {
  description = "Zona GCP para el cluster"
  type        = string
}

variable "network_id" {
  description = "ID de la red VPC"
  type        = string
}

variable "subnet_id" {
  description = "ID de la subred para el cluster"
  type        = string
}

variable "node_count" {
  description = "Número de nodos para el cluster"
  type        = number
}

variable "machine_type" {
  description = "Tipo de máquina para los nodos"
  type        = string
}

variable "service_account" {
  description = "Email de la cuenta de servicio para los nodos"
  type        = string
}