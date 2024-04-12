

data "terraform_remote_state" "network" {
  backend = var.network_backend_type
  config  = var.network_backend_config
}

resource "google_service_account" "default" {
  account_id   = local.compute_account_id
  display_name = format("%s-sa", local.compute_account_id)
}

resource "google_compute_instance" "vm" {
  name         = local.compute_name
  machine_type = "e2-micro"
  zone         = local.compute_zone
  labels = {
    managed_by = "terraform"
  }
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2004-lts"
    }
  }
  allow_stopping_for_update = true
  network_interface {
    network    = data.terraform_remote_state.network.outputs.network_self_link
    subnetwork = data.terraform_remote_state.network.outputs.subnets["${local.region}/${var.subnet_name}"].self_link
  }

  shielded_instance_config {
    enable_secure_boot = true
    enable_vtpm        = true
  }
}
