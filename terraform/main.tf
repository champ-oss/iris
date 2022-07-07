locals {
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}

module "lambda" {
  source                          = "github.com/champ-oss/terraform-aws-lambda.git?ref=v1.0.66-e3b8bd1"
  git                             = var.git
  name                            = "lambda"
  vpc_id                          = var.enable_vpc ? var.vpc_id : null
  private_subnet_ids              = var.enable_vpc ? var.private_subnet_ids : null
  enable_vpc                      = var.enable_vpc
  enable_function_url             = true
  function_url_authorization_type = "NONE"
  reserved_concurrent_executions  = var.reserved_concurrent_executions
  sync_image                      = true
  sync_source_repo                = "champtitles/iris"
  ecr_name                        = "${var.git}-lambda"
  ecr_tag                         = var.docker_tag
  tags                            = merge(local.tags, var.tags)
  environment = {
    ALLOWED_URLS = join(",", var.allowed_urls)
  }
}
