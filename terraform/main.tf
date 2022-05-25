locals {
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}

data "aws_region" "this" {}
data "aws_caller_identity" "this" {}

resource "random_string" "identifier" {
  length  = 5
  special = false
  upper   = false
  lower   = true
  number  = true
}

module "lambda" {
  depends_on                      = [null_resource.sync_dockerhub_ecr]
  source                          = "github.com/champ-oss/terraform-aws-lambda.git?ref=2a54a80fd659d856ce2ac7eff4390d0a9c2cdcbb"
  git                             = var.git
  name                            = random_string.identifier.result
  vpc_id                          = var.enable_vpc ? var.vpc_id : null
  private_subnet_ids              = var.enable_vpc ? var.private_subnet_ids : null
  enable_vpc                      = var.enable_vpc
  enable_function_url             = true
  function_url_authorization_type = "NONE"
  reserved_concurrent_executions  = var.reserved_concurrent_executions
  ecr_account                     = data.aws_caller_identity.this.account_id
  ecr_name                        = aws_ecr_repository.this.name
  ecr_tag                         = var.docker_tag
  tags                            = merge(local.tags, var.tags)
  environment = {
    ALLOWED_URLS = join(",", var.allowed_urls)
  }
}

resource "null_resource" "sync_dockerhub_ecr" {
  depends_on = [aws_ecr_repository.this]

  triggers = {
    ecr_name   = aws_ecr_repository.this.name
    docker_tag = var.docker_tag
  }

  provisioner "local-exec" {
    command     = "sh ${path.module}/sync_dockerhub_ecr.sh"
    interpreter = ["/bin/sh", "-c"]
    environment = {
      RETRIES     = 60
      SLEEP       = 10
      AWS_REGION  = data.aws_region.this.name
      SOURCE_REPO = "champtitles/iris"
      IMAGE_TAG   = var.docker_tag
      ECR_ACCOUNT = data.aws_caller_identity.this.account_id
      ECR_NAME    = aws_ecr_repository.this.name
    }
  }
}

resource "aws_ecr_repository" "this" {
  name = "${var.git}-${random_string.identifier.result}"
  tags = merge(local.tags, var.tags)

  image_scanning_configuration {
    scan_on_push = true
  }
}
