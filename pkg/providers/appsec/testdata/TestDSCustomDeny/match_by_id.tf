provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_custom_deny" "test" {
  config_id      = 43253
  custom_deny_id = "deny_custom_54994"
}

