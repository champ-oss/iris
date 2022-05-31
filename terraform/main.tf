locals {
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}

data "aws_region" "this" {}
data "aws_caller_identity" "this" {}

module "lambda" {
  depends_on                      = [module.ecr]
  source                          = "github.com/champ-oss/terraform-aws-lambda.git?ref=v1.0.19-1702466"
  git                             = var.git
  name                            = "lambda"
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

module "ecr" {
  source           = "github.com/champ-oss/terraform-aws-ecr.git?ref=v1.0.22-24fb4c0"
  name             = "${var.git}-lambda"
  sync_image       = true
  sync_source_repo = "champtitles/iris"
  sync_source_tag  = var.docker_tag
}
