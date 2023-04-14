provider "aws" {
  region = "us-east-1"
}

locals {
  git = "iris"
}

module "this" {
  source     = "../../"
  enable_vpc = false

  allowed_urls = [
    "about.google/google-in-america",
    "aws.amazon.com/console",
    "1.com/foo"
  ]
}