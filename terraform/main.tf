locals {
  dns = var.hostname != null ? "${var.hostname}.${var.domain}" : "${var.git}-${random_string.identifier.result}.${var.domain}"

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

module "acm" {
  source            = "github.com/champ-oss/terraform-aws-acm.git?ref=v1.0.17-66adf61"
  git               = var.git
  domain_name       = local.dns
  create_wildcard   = false
  zone_id           = var.zone_id
  enable_validation = true
}

module "alb" {
  source          = "github.com/champ-oss/terraform-aws-alb.git?ref=v1.0.4-e4392ee"
  git             = var.git
  certificate_arn = module.acm.arn
  subnet_ids      = var.public_subnet_ids
  vpc_id          = var.vpc_id
  internal        = false
  protect         = false
  tags            = merge(local.tags, var.tags)
}

module "lambda" {
  depends_on           = [null_resource.sync_dockerhub_ecr]
  source               = "github.com/champ-oss/terraform-aws-lambda.git?ref=v1.0.6-e2d9736"
  git                  = var.git
  name                 = random_string.identifier.result
  vpc_id               = var.vpc_id
  private_subnet_ids   = var.private_subnet_ids
  zone_id              = var.zone_id
  listener_arn         = module.alb.listener_arn
  lb_dns_name          = module.alb.dns_name
  lb_zone_id           = module.alb.zone_id
  enable_load_balancer = true
  enable_route53       = true
  enable_vpc           = true
  dns_name             = local.dns
  ecr_account          = data.aws_caller_identity.this.account_id
  ecr_name             = aws_ecr_repository.this.name
  ecr_tag              = var.docker_tag
  tags                 = merge(local.tags, var.tags)
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
      SLEEP       = 5
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
