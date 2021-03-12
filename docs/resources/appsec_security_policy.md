---
layout: "akamai"
page_title: "Akamai: SecurityPolicy"
subcategory: "Application Security"
description: |-
 SecurityPolicy
---

# akamai_appsec_security_policy

Use the `akamai_appsec_security_policy` resource to create a new security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to create a new security policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_security_policy" "security_policy_create" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  default_settings = var.default_settings
  security_policy_name = var.policy_name
  security_policy_prefix = var.policy_prefix
}

output "security_policy_create" {
  value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

* `version` - (Required) The version number of the configuration to use.

* `security_policy_name` - (Required) The name of the new security policy.

* `security_policy_prefix' - (Required) The four-character alphanumeric string prefix for the policy ID.

* `default_settings` - (Optional) Whether the new policy should use the default settings. If not supplied, defaults to true.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `security_policy_id` - The ID of the newly created security policy.

