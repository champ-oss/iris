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
  expected_header_key   = "X-TEST-HEADER"
  expected_header_value = "25920B896F8744F4A33D3262EC7DD3DE"
}

module "this" {
  source                = "../../"
  enable_vpc            = false
  tags                  = local.tags
  expected_header_key   = local.expected_header_key
  expected_header_value = local.expected_header_value
  allowed_urls = [
    "about.google/google-in-america",
    "aws.amazon.com/console",
    "1.com/foo"
  ]
}