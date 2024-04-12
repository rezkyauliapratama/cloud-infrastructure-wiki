module "vpc" {
  source       = "terraform-google-modules/network/google"
  version      = "9.0.0"
  project_id   = var.project_id
  network_name = "${var.network_name}-network"
  routing_mode = "REGIONAL"

  subnets = [
    {
      subnet_name   = "${var.network_name}-subnet-01"
      subnet_ip     = "10.10.0.0/24"
      subnet_region = var.region
      description   = "This subnet is managed by Terraform"
    },
    {
      subnet_name   = "${var.network_name}-subnet-02"
      subnet_ip     = "10.10.1.0/24"
      subnet_region = var.region
      description   = "This subnet is managed by Terraform"
    },
    {
      subnet_name   = "${var.network_name}-subnet-03"
      subnet_ip     = "10.10.2.0/24"
      subnet_region = var.region
      description   = "This subnet is managed by Terraform"
    }
  ]
}


resource "google_compute_firewall" "allow-ssh" {
  name    = "allow-ssh-firewall"
  network = module.vpc.network_name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"] #allow ssh without ips whitelist
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
