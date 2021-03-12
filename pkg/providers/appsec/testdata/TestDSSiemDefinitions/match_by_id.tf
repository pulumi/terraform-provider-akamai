provider "akamai" {
  edgerc = "~/.edgerc"
}



data "akamai_appsec_siem_definitions" "test" {
  siem_definition_name = "SIEM Version 01"
}


