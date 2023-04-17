output "function_url" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_url#function_url"
  value       = module.this.function_url
}

output "function_url_no_header" {
  description = "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function_url#function_url"
  value       = module.no_header.function_url
}

output "expected_header_key" {
  description = "Header key that must be on every request"
  value       = local.expected_header_key
}

output "expected_header_value" {
  description = "Header value that must be on every request"
  value       = local.expected_header_value
}