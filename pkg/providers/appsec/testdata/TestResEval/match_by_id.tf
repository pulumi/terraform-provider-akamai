provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_eval" "test" {
    config_id = 43253
    security_policy_id = "AAAA_81230"
    eval_operation = "START"
}

