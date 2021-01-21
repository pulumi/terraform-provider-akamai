---
layout: "akamai"
page_title: "Akamai: Eval Rule Condition / Exception"
subcategory: "Application Security"
description: |-
 Eval Rule Condition / Exception
---

# akamai_appsec_eval_rule_condition_exception

Use the `akamai_appsec_eval_rule_condition_exception` resource to create or modify an eval rule's conditions and exceptions.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to add condition-exception to an eval rule using a JSON input
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_eval_rule_condition_exception" "condition_exception" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  rule_id = var.rule_id
  condition_exception = file("${path.module}/condition_exception.json")
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rule_id` - (Required) The ID of the eval rule to use.

* `condition_exception` - (Required) The name of a file containing a JSON-formatted description of the conditions and exceptions to use ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalconditionsexceptions))

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

