provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_configuration" "test" {
  name = "Akamai Tools New"
  description = "TF Tools 1"
  contract_id= "C-1FRYVV3"
  group_id  = 64867
  host_names = ["rinaldi.sandbox.akamaideveloper.com",
        "sujala.sandbox.akamaideveloper.com"]
}

