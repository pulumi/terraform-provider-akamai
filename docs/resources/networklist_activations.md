---
layout: "akamai"
page_title: "Akamai: NetworkLists Activations"
subcategory: "Network Lists"
description: |-
 NetworkLists
---

# akamai_networklist_activations

Use the `akamai_networklist_activations` resource to activate a network list in either the STAGING or PRODUCTION
environment.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_networklist_network_lists" "network_lists_filter" {
  name = var.network_list
}

resource "akamai_networklist_activations" "activation" {
  network_list_id = data.akamai_networklist_network_lists.network_lists_filter.list[0]
  network = "STAGING"
  notes  = "TEST Notes"
  notification_emails = ["user@example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `network_list_id` - (Required) The ID of the network list to be activated

* `network` - (Optional) The network to be used, either `STAGING` or `PRODUCTION`. If not supplied, defaults to
  `STAGING`.

* `notes` - (Optional) A comment describing the activation.

* `notification_emails` - (Required) A bracketed, comma-separated list of email addresses that will be notified when the
  operation is complete.

## Attributes Reference

In addition to the arguments above, the following attribute is exported:

* `status` - The string `ACTIVATED` if the activation was successful, or a string identifying the reason why the network
  list was not activated.

