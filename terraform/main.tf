locals {
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = var.git
  }
}

module "hash" {
  source   = "github.com/champ-oss/terraform-git-hash.git?ref=v1.0.12-fc3bb87"
  path     = "${path.module}/.."
  fallback = ""
}

module "lambda" {
  source                          = "github.com/champ-oss/terraform-aws-lambda.git?ref=v1.0.115-77403a9"
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
  ecr_tag                         = module.hash.hash
  tags                            = merge(local.tags, var.tags)
  environment = {
    ALLOWED_URLS          = join(",", var.allowed_urls)
    EXPECTED_HEADER_KEY   = var.expected_header_key
    EXPECTED_HEADER_VALUE = var.expected_header_value
  }
}
