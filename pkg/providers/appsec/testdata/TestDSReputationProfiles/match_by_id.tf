provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_reputation_profiles" "test" {
  config_id = 43253
  version = 7
  reputation_profile_id = 12345
}