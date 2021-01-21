---
layout: "akamai"
page_title: "Akamai: Attack Group Action"
subcategory: "Application Security"
description: |-
 Attack Group Action
---

# akamai_appsec_attack_group_action

Use the `akamai_appsec_attack_group_action` resource to update what action should be taken when an attack group’s rule triggers. 

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to set the attack group action
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_attack_group_action" "attack_group_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  attack_group = var.attack_group
  attack_group_action = var.action
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `attack_group` - (Required) The ID of the attack group to use.

* `attack_group_action` - (Required) The action to be taken: `alert` to record the trigger of the event, `deny` to block the request, `deny_custom_{custom_deny_id}` to execute a custom deny action, or `none` to take no action.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

