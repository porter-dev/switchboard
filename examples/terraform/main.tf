variable "rds_username" {
    type = string
}

variable "rds_password" {
    type = string
}

variable "rds_host" {
    type = string
}

resource "null_resource" "example1" {
    triggers = {
    always_run = "${timestamp()}"
  }


  provisioner "local-exec" {
    command = "echo 'This is, like, an RDS example.'"
  }
}

output "rds_username" {
    value = var.rds_username
}

output "rds_password" {
    value = var.rds_password
}

output "rds_host" {
    value = var.rds_host
}