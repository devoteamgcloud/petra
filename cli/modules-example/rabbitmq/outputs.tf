output "dummy" {
  value = helm_release.master.metadata.0.name
}
