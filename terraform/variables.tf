variable "commit_sha" {
  description = "Git commit SHA of code expected to be deployed"
  type        = string
}

variable "allowed_urls" {
  description = "List of URLs that will be allowed to proxy"
  type        = list(string)
  default     = []
}

variable "tags" {
  description = "Map of tags to assign to resources"
  type        = map(string)
  default     = {}
}

variable "git" {
  description = "Name of the Git repo"
  type        = string
  default     = "iris"
}

variable "vpc_id" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group#vpc_id"
  type        = string
}

variable "public_subnet_ids" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb#subnets"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_cluster#subnet_ids"
  type        = list(string)
}

variable "hostname" {
  description = "Optional hostname for Iris. If omitted a random identifier will be used."
  type        = string
  default     = null
}

variable "domain" {
  description = "Route53 Domain"
  type        = string
}

variable "zone_id" {
  description = "Route53 Zone ID"
  type        = string
}

