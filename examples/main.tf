
provider "pcfom" {
  target_hostname = "${var.target_hostname}"
  token = "${var.token}"
  skip_ssl_validation = true
}

data "local_file" "director_config" {
  filename = "${path.module}/director-config.json"
}

resource "pcfom_director" "aws_cp" {
  director_config = "${data.local_file.director_config.content}"
}

variable "token" {
  type = "string"
}

variable "target_hostname" {
  type = "string"
}