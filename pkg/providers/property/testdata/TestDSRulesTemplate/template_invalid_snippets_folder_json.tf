provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_property_rules_template" "test" {
  template_file = "testdata/TestDSRulesTemplate/output/template_invalid_json.json"
}
