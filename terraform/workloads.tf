# ------------------------------------------------------------------------------
# 1. Kubernetes Service Accounts (KSAs)
# ------------------------------------------------------------------------------

resource "kubernetes_service_account_v1" "ksa_classifier" {
  metadata {
    name      = "ksa-classifier"
    namespace = "default"
    annotations = {
      "iam.gke.io/gcp-service-account" = google_service_account.sa_classifier.email
    }
  }
}

resource "kubernetes_service_account_v1" "ksa_auditor" {
  metadata {
    name      = "ksa-auditor"
    namespace = "default"
    annotations = {
      "iam.gke.io/gcp-service-account" = google_service_account.sa_auditor.email
    }
  }
}

resource "kubernetes_service_account_v1" "ksa_vuln" {
  metadata {
    name      = "ksa-vuln"
    namespace = "default"
    annotations = {
      "iam.gke.io/gcp-service-account" = google_service_account.sa_vuln.email
    }
  }
}

resource "kubernetes_service_account_v1" "ksa_reporter" {
  metadata {
    name      = "ksa-reporter"
    namespace = "default"
    annotations = {
      "iam.gke.io/gcp-service-account" = google_service_account.sa_reporter.email
    }
  }
}

# ------------------------------------------------------------------------------
# 2. SecretProviderClass (Connects CSI Driver to Secret Manager)
# ------------------------------------------------------------------------------
resource "kubectl_manifest" "secret_provider_class" {
  yaml_body = yamlencode({
    apiVersion = "secrets-store.csi.x-k8s.io/v1"
    kind       = "SecretProviderClass"
    metadata = {
      name      = "gemini-api-key-spc"
      namespace = "default"
    }
    spec = {
      provider = "gcp"
      parameters = {
        secrets = yamlencode([{
          resourceName = google_secret_manager_secret.gemini_api_key.secret_id
          fileName     = "key"
        }])
      }
      secretObjects = [{
        secretName = "gemini-api-key"
        type       = "Opaque"
        data = [{
          objectName = "key"
          key        = "key"
        }]
      }]
    }
  })
}

# ------------------------------------------------------------------------------
# 3. Deployments (Micro-Agent Architecture)
# ------------------------------------------------------------------------------

# --- Scope Classifier ---
resource "kubernetes_deployment_v1" "agent_classifier" {
  metadata { name = "agent-classifier" }
  spec {
    replicas = 1
    selector { match_labels = { app = "agent-classifier" } }
    template {
      metadata { labels = { app = "agent-classifier" } }
      spec {
        service_account_name = kubernetes_service_account_v1.ksa_classifier.metadata[0].name
        volume {
          name = "secrets-store-inline"
          csi {
            driver = "secrets-store.csi.k8s.io"
            read_only = true
            volume_attributes = { secretProviderClass = "gemini-api-key-spc" }
          }
        }
        container {
          image = var.image_repository
                    name  = "classifier"
                    # Hypothetical entrypoint command
                    command = ["/app/main"]
                    args    = ["--role=classifier", "--mode=server"]
                    
                    volume_mount {
                      name       = "secrets-store-inline"
                      mount_path = "/mnt/secrets-store"
                      read_only  = true
                    }
                    env {
                      name = "GEMINI_API_KEY"
                      value_from {
                        secret_key_ref {
                          name = "gemini-api-key"
                          key  = "key"
                        }
                      }
                    }
                    env {
                      name  = "PROJECT_ID"
                      value = var.project_id
                    }
                    
                    port { container_port = 8080 }
                    
                    resources {
                      limits = {
                        cpu    = "500m"
                        memory = "512Mi"
                      }
                      requests = {
                        cpu    = "250m"
                        memory = "256Mi"
                      }
                    }
                    security_context {
                      run_as_non_root = true
                      allow_privilege_escalation = false
                      capabilities {
                        drop = ["ALL"]
                      }
                    }
          liveness_probe {
            http_get {
              path = "/healthz"
              port = 8080
            }
            initial_delay_seconds = 3
            period_seconds        = 3
          }
          image_pull_policy = "Always"
        }
      }
    }
  }
}

# --- Regulatory Auditor ---
resource "kubernetes_deployment_v1" "agent_auditor" {
  metadata { name = "agent-auditor" }
  spec {
    replicas = 1
    selector { match_labels = { app = "agent-auditor" } }
    template {
      metadata { labels = { app = "agent-auditor" } }
      spec {
        service_account_name = kubernetes_service_account_v1.ksa_auditor.metadata[0].name
        volume {
          name = "secrets-store-inline"
          csi {
            driver = "secrets-store.csi.k8s.io"
            read_only = true
            volume_attributes = { secretProviderClass = "gemini-api-key-spc" }
          }
        }
        container {
          image = var.image_repository
                                        name  = "auditor"
                                        command = ["/app/main"]
                                        args    = ["--role=auditor", "--mode=server"]
                    
                                        volume_mount {
                                          name       = "secrets-store-inline"
                                          mount_path = "/mnt/secrets-store"
                                          read_only  = true
                                        }
                                        env {
                                          name = "GEMINI_API_KEY"
                                          value_from {
                                            secret_key_ref {
                                              name = "gemini-api-key"
                                              key  = "key"
                                            }
                                          }
                                        }
                                        env {
                                          name  = "PROJECT_ID"
                                          value = var.project_id
                                        }
                              port { container_port = 8080 }
                    
                              resources {
                                limits = {
                                  cpu    = "500m"
                                  memory = "512Mi"
                                }
                                requests = {
                                  cpu    = "250m"
                                  memory = "256Mi"
                                }
                              }
          security_context {
            run_as_non_root = true
            allow_privilege_escalation = false
            capabilities {
              drop = ["ALL"]
            }
          }
          liveness_probe {
            http_get {
              path = "/healthz"
              port = 8080
            }
            initial_delay_seconds = 3
            period_seconds        = 3
          }
          image_pull_policy = "Always"
        }
      }
    }
  }
}

# --- Vulnerability Watchdog ---
resource "kubernetes_deployment_v1" "agent_vuln" {
  metadata { name = "agent-vuln" }
  spec {
    replicas = 1
    selector { match_labels = { app = "agent-vuln" } }
    template {
      metadata { labels = { app = "agent-vuln" } }
      spec {
        service_account_name = kubernetes_service_account_v1.ksa_vuln.metadata[0].name
        volume {
          name = "secrets-store-inline"
          csi {
            driver = "secrets-store.csi.k8s.io"
            read_only = true
            volume_attributes = { secretProviderClass = "gemini-api-key-spc" }
          }
        }
        container {
          image = var.image_repository
          name  = "vuln"
          command = ["/app/main"] 
          args    = ["--role=vuln", "--mode=server"]

          volume_mount {
            name       = "secrets-store-inline"
            mount_path = "/mnt/secrets-store"
            read_only  = true
          }
          env {
            name = "GEMINI_API_KEY"
            value_from {
              secret_key_ref {
                name = "gemini-api-key"
                key  = "key"
              }
            }
          }
          env {
            name  = "PROJECT_ID"
            value = var.project_id
          }

          port { container_port = 8080 }

          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }
          security_context {
            run_as_non_root = true
            allow_privilege_escalation = false
            capabilities {
              drop = ["ALL"]
            }
          }
          liveness_probe {
            http_get {
              path = "/healthz"
              port = 8080
            }
            initial_delay_seconds = 3
            period_seconds        = 3
          }
          image_pull_policy = "Always"
        }
      }
    }
  }
}

# --- Compliance Reporter ---
resource "kubernetes_deployment_v1" "agent_reporter" {
  metadata { name = "agent-reporter" }
  spec {
    replicas = 1
    selector { match_labels = { app = "agent-reporter" } }
    template {
      metadata { labels = { app = "agent-reporter" } }
      spec {
        service_account_name = kubernetes_service_account_v1.ksa_reporter.metadata[0].name
        volume {
          name = "secrets-store-inline"
          csi {
            driver = "secrets-store.csi.k8s.io"
            read_only = true
            volume_attributes = { secretProviderClass = "gemini-api-key-spc" }
          }
        }
        container {
          image = var.image_repository
          name  = "reporter"
          command = ["/app/main"] 
          args    = ["--role=reporter", "--mode=server"]

          volume_mount {
            name       = "secrets-store-inline"
            mount_path = "/mnt/secrets-store"
            read_only  = true
          }
          env {
            name = "GEMINI_API_KEY"
            value_from {
              secret_key_ref {
                name = "gemini-api-key"
                key  = "key"
              }
            }
          }
          env {
            name  = "PROJECT_ID"
            value = var.project_id
          }

          port { container_port = 8080 }

          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }
          security_context {
            run_as_non_root = true
            allow_privilege_escalation = false
            capabilities {
              drop = ["ALL"]
            }
          }
          liveness_probe {
            http_get {
              path = "/healthz"
              port = 8080
            }
            initial_delay_seconds = 3
            period_seconds        = 3
          }
          image_pull_policy = "Always"
        }
      }
    }
  }
}

# ------------------------------------------------------------------------------
# 4. Services (Internal Communication)
# ------------------------------------------------------------------------------

resource "kubernetes_service_v1" "svc_classifier" {
  metadata { name = "svc-classifier" }
  spec {
    selector = { app = "agent-classifier" }
    port {
      port        = 80
      target_port = 8080
    }
    type = "ClusterIP"
  }
}

resource "kubernetes_service_v1" "svc_auditor" {
  metadata { name = "svc-auditor" }
  spec {
    selector = { app = "agent-auditor" }
    port {
      port        = 80
      target_port = 8080
    }
    type = "ClusterIP"
  }
}

resource "kubernetes_service_v1" "svc_vuln" {
  metadata { name = "svc-vuln" }
  spec {
    port {
      port        = 80
      target_port = 8080
    }
    type = "ClusterIP"
  }
}

resource "kubernetes_service_v1" "svc_reporter" {
  metadata { name = "svc-reporter" }
  spec {
    selector = { app = "agent-reporter" }
    port {
      port        = 80
      target_port = 8080
    }
    type = "ClusterIP"
  }
}