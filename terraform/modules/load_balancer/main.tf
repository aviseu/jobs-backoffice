resource "google_compute_region_network_endpoint_group" "group" {
  for_each              = var.backends

  name                  = "${each.key}-neg"
  project               = var.project_id
  region                = var.region
  network_endpoint_type = "SERVERLESS"

  cloud_run {
    service = each.value
  }
}

resource "google_compute_region_backend_service" "backend" {
  depends_on            = [google_compute_region_network_endpoint_group.group]
  for_each              = var.backends

  region                = var.region
  project               = var.project_id
  name                  = "${each.key}-backend"
  protocol              = "HTTPS"
  port_name             = "http"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  timeout_sec           = 30

  backend {
    group               = google_compute_region_network_endpoint_group.group[each.key].self_link
    balancing_mode      = "UTILIZATION"
    capacity_scaler     = 1.0
  }

  connection_draining_timeout_sec = 0
}

data "google_compute_address" "static_ip" {
  name   = var.address_name
  region = var.region
}

resource "google_compute_region_url_map" "load_balancer" {
  depends_on = [google_compute_region_backend_service.backend]

  name    = var.load_balancer_name
  project = var.project_id
  region  = var.region

  default_service = google_compute_region_backend_service.backend[var.default_backend].self_link

  dynamic "host_rule" {
    for_each = var.routes

    content {
      hosts           = [host_rule.value.domain]
      path_matcher    = "path-matcher-${host_rule.key}"
    }
  }

  dynamic "path_matcher" {
    for_each = var.routes

    content {
      name            = "path-matcher-${path_matcher.key}"
      default_service = google_compute_region_backend_service.backend[var.default_backend].self_link

      dynamic "path_rule" {
        for_each = path_matcher.value.paths

        content {
          paths   = [path_rule.key]
          service = google_compute_region_backend_service.backend[path_rule.value].self_link
        }
      }
    }
  }
}

data "google_compute_region_ssl_certificate" "certificate" {
  for_each = var.routes

  project = var.project_id
  region = var.region
  name = each.value.certificate
}

resource "google_compute_region_target_https_proxy" "proxy" {
  depends_on = [google_compute_region_url_map.load_balancer, data.google_compute_region_ssl_certificate.certificate]
  for_each = data.google_compute_region_ssl_certificate.certificate

  project           = var.project_id
  region            = var.region
  name              = "${var.load_balancer_name}-${each.key}-proxy"
  url_map           = google_compute_region_url_map.load_balancer.id
  ssl_certificates  = [data.google_compute_region_ssl_certificate.certificate[each.key].id]
}

resource "google_compute_forwarding_rule" "google_compute_forwarding_rule" {
  for_each = var.routes

  name                  = "${var.load_balancer_name}-forwarding-rule"
  region                = var.region
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  port_range            = "443-443"
  target                = google_compute_region_target_https_proxy.proxy[each.key].self_link
  network               = "default"
  network_tier          = "STANDARD"
  ip_address            = data.google_compute_address.static_ip.address
}
