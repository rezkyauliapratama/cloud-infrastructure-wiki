
variable "network_backend_type" {
  type    = string
  default = "gcs"
}

variable "network_backend_config" {
  # type        = object({bucket = string, prefix = string, encryption_key = string})
  type = map(string)
}


variable "project_id" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "region" {
  type = string
}

variable "preemptible" {
  type        = bool
  description = "Allow the instance to be preempted"
  default     = true
}


variable "kubernetes_version" {
  type        = string
  default     = "latest"
}

variable "disk_size" {
    type = number
    default = 50
}

variable "cluster_regional" {
    type = bool
    default = true
}

variable "default_node_machine_type" {
    type = string
    default = "e2-standard-2"
}

variable "default_node_min_count" {
    type = number
    default = 1
}

variable "default_node_max_count" {
    type = number
    default = 5
}

variable "subnet_name" {
    type = string
}

variable "subnet_secondary_pod" {
    type = string
}

variable "subnet_secondary_svc" {
    type = string
}

variable "zone" {
    type = string
    default = ""
}

variable "spot" {
    type = bool
    default = false
}