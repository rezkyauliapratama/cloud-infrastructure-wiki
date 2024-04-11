

module "vpc" {
  source       = "terraform-google-modules/network/google"
  version      = "9.0.0"
  project_id   = var.project_id
  network_name = "${var.network_name}-network"
  routing_mode = "REGIONAL"

  subnets = [
    {
      subnet_name   = "${var.network_name}-subnet"
      subnet_ip     = local.subnet_ip
      subnet_region = var.region
      description   = "This subnet is managed by Terraform"
    }
  ]
  secondary_ranges = {
    "${var.network_name}-subnet" = [
      {
        range_name    = "ip-range-pods"
        ip_cidr_range = local.ip_range_pods
      },
      {
        range_name    = "ip-range-svc"
        ip_cidr_range = local.ip_range_svc
      },
    ]
  }
}

module "cloud-nat" {
  source        = "terraform-google-modules/cloud-nat/google"
  version       = "5.0.0"
  project_id    = var.project_id
  region        = var.region
  router        = "${var.network_name}-router"
  network       = module.vpc.network_self_link
  create_router = true
}

