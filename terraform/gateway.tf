# Gateway Resource
resource "kubectl_manifest" "gateway" {
  yaml_body = yamlencode({
    apiVersion = "gateway.networking.k8s.io/v1beta1"
    kind       = "Gateway"
    metadata = {
      name      = "external-gateway"
      namespace = "default"
    }
    spec = {
      gatewayClassName = "gke-l7-global-external-managed"
      listeners = [{
        name     = "http"
        protocol = "HTTP"
        port     = 80
        allowedRoutes = {
          namespaces = {
            from = "Same"
          }
        }
      }]
    }
  })
}

# HTTPRoute Resource attaching the Service + Security Policy
resource "kubectl_manifest" "http_route" {
  yaml_body = yamlencode({
    apiVersion = "gateway.networking.k8s.io/v1beta1"
    kind       = "HTTPRoute"
    metadata = {
      name      = "agent-route"
      namespace = "default"
    }
    spec = {
      parentRefs = [{
        name = "external-gateway"
      }]
      rules = [{
        matches = [{
          path = {
            type  = "PathPrefix"
            value = "/"
          }
        }]
        backendRefs = [{
          name = "agent-cra-service"
          port = 80
        }]
      }]
    }
  })
}

# GCPBackendPolicy to attach Cloud Armor to the Backend Service created by Gateway
resource "kubectl_manifest" "backend_policy" {
  yaml_body = yamlencode({
    apiVersion = "networking.gke.io/v1"
    kind       = "GCPBackendPolicy"
    metadata = {
      name      = "agent-backend-policy"
      namespace = "default"
    }
    spec = {
      default = {
        securityPolicy = google_compute_security_policy.agent_armor.name
      }
      targetRef = {
        group = ""
        kind  = "Service"
        name  = "agent-cra-service"
      }
    }
  })
}
