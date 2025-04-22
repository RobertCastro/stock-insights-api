variable "network_name" {
  description = "Nombre de la red VPC"
  type        = string
}

variable "region" {
  description = "Región GCP para recursos de red"
  type        = string
}

variable "private_subnet_cidr" {
  description = "CIDR para subred privada"
  type        = string
}