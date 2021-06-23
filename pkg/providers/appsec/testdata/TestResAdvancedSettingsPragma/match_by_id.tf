provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_advanced_settings_pragma_header" "test" {
  config_id = 43253
  pragma_header = <<-EOF
{"action":"REMOVE"}
EOF
}



