
provider "pcfom" {
  target_hostname = "${var.target_hostname}"
  token = "${var.token}"
  skip_ssl_validation = true
}

data "local_file" "director_config" {
  filename = "${path.module}/director-config.json"
}

data "local_file" "minio_tile_config" {
  filename = "${path.module}/minio-tile-config.json"
}

resource "pcfom_director" "aws_cp" {
  director_config = "${data.local_file.director_config.content}"
}

resource "pcfom_tile" "minio" {
  product_name = "minio"
  tile_file = "/full/path/to/minio-internal-blobstore-1.0.4.pivotal"
  tile_config = "${data.local_file.minio_tile_config.content}"
}

variable "token" {
  type = "string"
}

variable "target_hostname" {
  type = "string"
}
