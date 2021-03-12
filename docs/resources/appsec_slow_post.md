---
layout: "akamai"
page_title: "Akamai: Slow Post"
subcategory: "Application Security"
description: |-
 Slow Post
---

# akamai_appsec_slow_post

Use the `akamai_appsec_slow_post` data source to update the slow post protection settings for a given security configuration version and policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user would like to set the slow post protection settings for a given security configuration and version
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_slow_post" "slow_post" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  slow_rate_action = "alert"
  slow_rate_threshold_rate = 10
  slow_rate_threshold_period = 30
  duration_threshold_timeout = 20
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `slow_rate_action` - (Required) The action that the rule should trigger (either `alert` or `abort`).

* `slow_rate_threshold_rate` - (Required) The average rate in bytes per second over the period specified by `period` before the specified `action` is triggered.

* `slow_rate_threshold_period` - (Required) The slow rate period value: the amount of time in seconds that the server should accept a request to determine whether a POST request is too slow. 

* `duration_threshold_timeout` - (Required) The time in seconds before the first eight kilobytes of the POST body must be received to avoid triggering the specified `action`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

