data "cloudavenue_jobs" "example" {
  id = "fb064495-457d-40d4-8e53-79fe3824ca96"
}

output "jobs" {
  value = data.cloudavenue_jobs.example
}
