output "service_name" {
  value = google_cloud_run_v2_service.service.name
}

output "backend" {
  value = google_compute_region_backend_service.backend.self_link
}
