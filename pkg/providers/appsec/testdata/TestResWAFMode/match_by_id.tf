provider "akamai" {
  edgerc = "~/.edgerc"
}



resource "akamai_appsec_waf_mode" "test" {
    config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
    mode = "AAG"
}

output "configsedge_post_output_text" {
  value = akamai_appsec_waf_mode.test.output_text
}
