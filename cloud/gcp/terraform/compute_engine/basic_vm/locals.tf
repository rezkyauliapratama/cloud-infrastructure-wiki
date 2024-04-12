locals {
  compute_name       = "basic-vm"
  compute_zone       = format("%s-a", local.region)
  project_id         = var.project_id
  region             = "asia-southeast2"
  compute_account_id = "compute-operational"
}
