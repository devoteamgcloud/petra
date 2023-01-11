resource "google_storage_bucket" "auto_expire" {
  project       = var.project_id
  name          = var.name
  location      = var.location
  force_destroy = true

  public_access_prevention = "enforced"
}
