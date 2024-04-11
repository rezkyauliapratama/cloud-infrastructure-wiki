module "enabled_google_apis" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"

  project_id                  = var.project_id
  disable_services_on_destroy = false

  activate_apis = [
    "iam.googleapis.com",
    "compute.googleapis.com",
    "containerregistry.googleapis.com",
    "container.googleapis.com", # mandatory for GKE 
    "binaryauthorization.googleapis.com",
    "iap.googleapis.com",
    "servicenetworking.googleapis.com"
  ]
}
