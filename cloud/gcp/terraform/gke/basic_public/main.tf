data "terraform_remote_state" "network" {
  backend = var.network_backend_type
  config  = var.network_backend_config
}

locals {
  secondary_ranges = data.terraform_remote_state.network.outputs.subnets["${var.region}/${var.subnet_name}"].secondary_ip_range
}

# google_client_config and kubernetes provider must be explicitly specified like the following.
data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${module.gke.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(module.gke.ca_certificate)
}

# reference : https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/blob/master/examples/simple_regional/main.tf
module "gke" {
  source     = "terraform-google-modules/kubernetes-engine/google"
  version    = "27.0.0"
  project_id = module.enabled_google_apis.project_id
  name       = var.cluster_name

  region             = var.region
  regional           = var.cluster_regional
  network            = data.terraform_remote_state.network.outputs.network_name
  network_project_id = ""
  zones              = var.cluster_regional ? [] : ["${var.region}-a"]



  kubernetes_version = var.kubernetes_version

  release_channel = "STABLE"


  subnetwork        = var.subnet_name
  ip_range_pods     = local.secondary_ranges[index(local.secondary_ranges.*.range_name, var.subnet_secondary_pod)].range_name
  ip_range_services = local.secondary_ranges[index(local.secondary_ranges.*.range_name, var.subnet_secondary_svc)].range_name


  add_cluster_firewall_rules = true
  firewall_priority          = 1000

  horizontal_pod_autoscaling = false
  http_load_balancing        = false
  filestore_csi_driver       = false
  gce_pd_csi_driver          = true
  network_policy             = true

  maintenance_start_time = "18:00"

  remove_default_node_pool          = true
  disable_legacy_metadata_endpoints = true

  cluster_autoscaling = {
    enabled             = false # disable cluster autoscalling to prevent auto-provisioning node that sometimes can be tricky to control the cost and nodepool in your cluster.
    # all below items in this object is mandatory for cluster_autoscaling. So we need to define it regardless you enable or not cluster_autoscaling.
    autoscaling_profile = "OPTIMIZE_UTILIZATION"
    max_cpu_cores       = 80
    min_cpu_cores       = 0
    max_memory_gb       = 80
    min_memory_gb       = 0
    gpu_resources       = []
    auto_repair   = true
    auto_upgrade  = true
  }


  node_pools = [
    {
      name          = "default-node-pool"
      min_count     = var.default_node_min_count
      max_count     = var.default_node_max_count
      auto_upgrade  = true
      node_metadata = "GKE_METADATA"

      machine_type    = var.default_node_machine_type
      local_ssd_count = 0
      disk_size_gb    = var.disk_size
      disk_type       = "pd-ssd"
      image_type      = "COS_CONTAINERD"
      enable_gcfs     = true
      enable_gvnic    = true
      auto_repair     = true
      auto_upgrade    = true
      preemptible     = var.preemptible
      spot            = var.spot
      autoscaling     = true
    }
  ]

  node_pools_labels = {
    all = {
      default-node-pool = true
      cluster           = var.cluster_name
      managed_by        = "terraform"
    }

    default-node-pool = {
      default-node-pool = true
      cluster           = var.cluster_name
      managed_by        = "terraform"
    }
  }

  node_pools_tags = {
    default-node-pool = [
      "default-node-pool-${var.cluster_name}"
    ]
  }

  node_pools_oauth_scopes = {
    all = [
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/servicecontrol",
    ]

    default-node-pool = [
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/servicecontrol",
    ]
  }


  stub_domains         = {}
  upstream_nameservers = []

  logging_service    = "none" # No logs sent to Cloud Logging; no log collection agent installed in the cluster.
  monitoring_service = "none" # No metrics sent to Cloud Monitoring; no metric collection agent installed in the cluster.
  monitoring_enable_managed_prometheus = false
  create_service_account = true
  grant_registry_access  = true


  issue_client_certificate = false

  cluster_resource_labels = {
    cluster    = var.cluster_name
    managed_by = "terraform"
  }

  default_max_pods_per_node = 110

  database_encryption = [
    {
      state    = "DECRYPTED"
      key_name = ""
    }
  ]

  enable_binary_authorization = true

  enable_vertical_pod_autoscaling = false

  identity_namespace = "${var.project_id}.svc.id.goog"

  enable_shielded_nodes = true
}
