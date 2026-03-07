# Cloud Armor Security Policy (Model Armor Implementation)
resource "google_compute_security_policy" "agent_armor" {
  name = "agent-armor-policy"

  # Rule 1: Allow specific traffic (Placeholder for complex AI protection rules)
  rule {
    action   = "allow"
    priority = "1000"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "Allow access"
  }

  # Rule 2: SQL Injection Protection (Standard WAF)
  rule {
    action   = "deny(403)"
    priority = "900"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('sqli-v33-stable')"
      }
    }
    description = "Block SQL Injection"
  }

  # Rule 3: Cross-Site Scripting Protection (Standard WAF)
  rule {
    action   = "deny(403)"
    priority = "901"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('xss-v33-stable')"
      }
    }
    description = "Block XSS"
  }
}
