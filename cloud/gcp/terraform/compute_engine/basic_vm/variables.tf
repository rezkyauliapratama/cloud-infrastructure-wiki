
variable "network_backend_type" {
  type    = string
  default = "gcs"
}

variable "network_backend_config" {
  # type        = object({bucket = string, prefix = string, encryption_key = string})
  type = map(string)
}

variable "zone" {
  type    = string
  default = "a"
}

variable "subnet_name" {
    type = string
}

variable "image_project" {
  type    = string
  default = "debian-cloud"
}

variable "image_family" {
  type    = string
  default = "debian-10"
}

variable "machine_type" {
  type    = string
  default = "e2-micro"
}

variable "bastion_members" {
  type = list(string)

  description = "List of IAM resources to allow access to the bastion host"
  default     = []
}

variable "preemptible" {
  type        = bool
  description = "Allow the instance to be preempted"
  default     = false
}

variable "project_id" {
  type = string
}