variable "namespace" {
  type = string
}
variable "name" {
  type    = string
  default = "rabbitmq"
}
variable "name_suffix" {
  type    = string
  default = ""
}
variable "hosts" {
  type = string
}
variable "annotations" {
  type    = string
  default = ""
}
variable "cluster" {
  type = string
}
variable "chart_version" {
  type    = string
  default = null
}
variable "helm3" {
  type = object({
    checksumSecretOverride = string
    releaseServiceOverride = string
  })
}

variable "optionalValuesFile" {
  type        = string
  default     = null
  description = "Use additional yaml value file"
}

variable "prometheus_name" {
  type = string
  default = "operated"
}
