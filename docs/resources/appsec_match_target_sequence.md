---
layout: "akamai"
page_title: "Akamai: MatchTargetSequence"
subcategory: "Application Security"
description: |-
  MatchTargetSequence
---

# akamai_appsec_match_target_sequence


The `akamai_appsec_match_target_sequence` resource allows you to specify the order in which match targets are applied within a given security configuration.


## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_match_target_sequence" "match_target_sequence" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  match_target_sequence =  file("${path.module}/match_targets.json")
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `match_target_sequence` - (Required) The name of a JSON file containing the sequence of all match targets defined for the specified security configuration ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putsequence)).

## Attribute Reference

In addition to the arguments above, the following attribute is exported:

* None

