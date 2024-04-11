
output "router_name" {
  description = "Name of the router that was created"
  value       = module.cloud-nat.router_name
}

output "network_name" {
  value       = module.vpc.network_name
  description = "The name of the VPC being created"
}

output "subnets" {
  value       = module.vpc.subnets
  description = "The name of the VPC subnet being created"
}


output "vpc_network_uri" {
  description = "Location of vpc network"
  value       = module.vpc.network_self_link
}

output "subnets_secondary_ranges" {
  description = "Location of private service network"
  value       = module.vpc.subnets_secondary_ranges
}

output "project_id" {
  value = var.project_id
}

