variable "project_id" {
    type = string
}

variable "name" {
    type = string
    description = "Bucket Name"
    default = "test-bucket-module"
}

variable "location" {
    type = string
    description = "Bucket Location"
    default = "EU"
}
