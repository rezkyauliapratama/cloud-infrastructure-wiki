output "cluster_name" {
  description = "Cluster name"
  value       = module.gke.name
}

output "location" {
  description = "Cluster location (region if regional cluster, zone if zonal cluster)"
  value       = module.gke.location
}

output "type" {
  description = "Cluster Type"
  value       = module.gke.type
}

output "k8s_master_version" {
  description = "Kubernetes master version"
  value       = module.gke.master_version
}

output "cluster_node_pools" {
  description = "Kubernetes cluster node pools"
  value       = module.gke.node_pools_names

}

output "endpoint" {
  sensitive   = true
  description = "Cluster endpoint"
  value       = module.gke.endpoint
}

output "master_authorized_networks_config" {
  description = "Networks from which access to master is permitted"
  value       = module.gke.master_authorized_networks_config
}

output "ca_certificate" {
  sensitive   = true
  description = "Cluster ca certificate (base64 encoded)"
  value       = module.gke.ca_certificate
}

output "get_credentials_command" {
  description = "gcloud get-credentials command to generate kubeconfig for the private cluster"
  value       = format("gcloud container clusters get-credentials --project %s --zone %s --internal-ip %s", var.project_id, module.gke.location, module.gke.name)
}
