provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_networklist_network_lists" "test" {
   name = "40996_ARTYLABWHITELIST"
   type = "IP"
}
