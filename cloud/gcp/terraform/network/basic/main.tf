module "vpc" {
  source       = "terraform-google-modules/network/google"
  version      = "9.0.0"
  project_id   = var.project_id
  network_name = "${var.network_name}-network"
  routing_mode = "REGIONAL"

  subnets = [
    {
      subnet_name           = "${var.network_name}-subnet-01"
      subnet_ip             = "10.10.0.0/24"
      subnet_region         = var.region
      description           = "This subnet is managed by Terraform"
    },
    {
      subnet_name           = "${var.network_name}-subnet-02"
      subnet_ip             = "10.10.1.0/24"
      subnet_region         = var.region
      description           = "This subnet is managed by Terraform"
    },
    {
      subnet_name           = "${var.network_name}-subnet-03"
      subnet_ip             = "10.10.2.0/24"
      subnet_region         = var.region
      description           = "This subnet is managed by Terraform"
    }
  ]
}
