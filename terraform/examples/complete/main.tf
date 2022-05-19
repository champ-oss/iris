provider "aws" {
  region = "us-east-1"
}

data "aws_route53_zone" "this" {
  name = "oss.champtest.net."
}

locals {
  git = "iris"
}

module "vpc" {
  source                   = "github.com/champ-oss/terraform-aws-vpc.git?ref=v1.0.9-ca0a300"
  git                      = local.git
  availability_zones_count = 2
  retention_in_days        = 1
  create_private_subnets   = false
}

module "this" {
  source             = "../../"
  docker_tag         = var.docker_tag
  domain             = data.aws_route53_zone.this.name
  private_subnet_ids = module.vpc.private_subnets_ids
  public_subnet_ids  = module.vpc.private_subnets_ids
  vpc_id             = module.vpc.vpc_id
  zone_id            = data.aws_route53_zone.this.zone_id

  allowed_urls = [
    "about.google/google-in-america",
    "aws.amazon.com/console",
    "1.com/foo"
  ]
}