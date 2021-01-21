provider "akamai" {
  edgerc = "~/.edgerc"
}

resource  "akamai_appsec_ip_geo" "test" {
  config_id = 43253
  version = 7
  security_policy_id =  "AAAA_81230"
  mode = "block"
  geo_network_lists= ["40731_BMROLLOUTGEO","44831_ECSCGEOBLACKLIST"]
  ip_network_lists= ["49181_ADTIPBLACKLIST","49185_ADTWAFBYPASSLIST"]
  exception_ip_network_lists= ["68762_ADYEN","69601_ADYENPRODWHITELIST"]
}
//allow


