resource "random_id" "default" {
  byte_length = 4
}

locals {
  chart_version = coalesce(var.chart_version, "7.7.0")
}

resource "helm_release" "master" {
  repository = "${path.module}/charts/${local.chart_version}"
  chart      = "rabbitmq"
  name       = var.name
  namespace  = var.namespace
  timeout    = 600
  version    = local.chart_version

  values = [
    file("${path.module}/charts/${local.chart_version}/rabbitmq/values-production.yaml"),
    file("${path.module}/configurations/default.yaml"),
    templatefile(var.optionalValuesFile, {prometheus_name = var.prometheus_name})
  ]
  set {
    name  = "annotations"
    value = var.annotations
  }
  set {
    name  = "fullnameOverride"
    value = "${var.name}-${coalesce(var.name_suffix, random_id.default.hex)}"
  }
  set {
    name  = "test"
    value = "test2"
  }
  set {
    name  = "checksumSecretOverride" #TODO: REMOVE ON NEXT RELEASE THAT NEED A K8S REPLACE
    value = var.helm3 == null ? "" : var.helm3.checksumSecretOverride
  }
  set {
    name  = "releaseServiceOverride" #TODO: REMOVE ON NEXT RELEASE THAT NEED A K8S REPLACE
    value = var.helm3 == null ? "" : var.helm3.releaseServiceOverride
  }
}
