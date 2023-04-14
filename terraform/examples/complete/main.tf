terraform {
  backend "s3" {}
}

provider "aws" {
  region = "us-east-2"
}

locals {
  git = "iris"
  tags = {
    cost    = "shared"
    creator = "terraform"
    git     = local.git
  }
}

module "this" {
  source     = "../../"
  enable_vpc = false
  tags       = local.tags
  allowed_urls = [
    "about.google/google-in-america",
    "aws.amazon.com/console",
    "1.com/foo"
  ]
}