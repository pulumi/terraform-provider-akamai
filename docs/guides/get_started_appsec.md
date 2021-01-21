---
layout: "akamai"
page_title: "Akamai: Get Started with Application Security"
description: |-
  Get Started with Akamai Application Security using Terraform
---

# Get Started with Application Security

The Akamai Provider for Terraform provides you the ability to automate the creation, deployment, and management of security configurations, custom rules, match targets and other application security resources.

To get more information about Application Security, see the [API documentation](https://developer.akamai.com/api/cloud_security/application_security/v1.html)

## Configure the Terraform Provider

Set up your .edgerc credential files as described in [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started), and include read-write permissions for the Application Security API. 

1. Create a new folder called `terraform`
1. Inside the new folder, create a new file called `akamai.tf`.
1. Add the provider configuration to your `akamai.tf` file:

```hcl
provider "akamai" {
	edgerc = "~/.edgerc"
	config_section = "appsec"
}
```

## Prerequisites

To manage Application Security resources, you need to obtain at a minimum the following information:

* **Configuration ID**: The ID of the specific security configuration under which the resources are defined.

For certain resources, you will also need other information, such as the version number of the security configuration. The process of obtaining this information is described below.

## Retrieving Security Configuration Information

You can obtain the name and ID of the existing security configurations using the [`akamai_appsec_configuration`](../data-sources/appsec_configuration.md) data source. This data source can be used with no additional parameters to output information about all security configurations associated with your account. Add the following to your `akamai.tf` file:

```hcl
data "akamai_appsec_configuration" "configurations" {
}

output "configuration_list" {
  value = data.akamai_appsec_configuration.configurations.output_text
}
```

Once you have saved the file, switch to the terminal and initialize Terraform using the command:

```bash
$ terraform init
```

This command will install the latest version of the Akamai provider, as well as any other providers necessary. To update the Akamai provider version after a new release, simply run `terraform init` again.

## Test Your Configuration

To test your configuration, use `terraform plan`:

```bash
$ terraform plan
```

This command will make Terraform create a plan for the work it will do based on the configuration file. This will not actually make any changes and is safe to run as many times as you like.

## Apply Changes

To actually display the configuration information, or to create or modify resources as described further in this guide, we need to instruct Terraform to `apply` the changes outlined in the plan. To do this, in the terminal, run the command:

```bash
$ terraform apply
```

Once this command has been executed, Terraform will display to the terminal window a formatted list of all existing security configurations under your account, including for each its name and ID (`config_id`), the number of its most recently created version, and the number of the version currently active in staging and production, if applicable.

When you have identified the desired security configuration by name, you can load that specific configuration into Terraform's state. To do this, edit your `akamai.tf` file to add the `name` parameter to the `akamai_appsec_configuration` data block using the desired configuration name as its value, and change the `output` block so that it gives just the `config_id` attribute of the configuration. After these changes, the portion of your file below the initial `provider` block will look like this:

```hcl
data "akamai_appsec_configuration" "configuration" {
  name = "Example"
}

output "ID" {
  value = data.akamai_appsec_configuration.configuration.config_id
}
```

If you run `terraform apply` on this file, you should see the `config_id` value of the specific configuration displayed on your terminal.

## Displaying Information About a Specific Configuration

The provider's [`akamai_appsec_export_configuration`](../data-sources/appsec_export_configuration.md) data source can diplay complete information about a specific configuration, including attributes such as custom rules, selected hostnames, etc. To show these two types of data for the most recent version of your selected configuration, add the following blocks to your `akamai.tf` file:

```hcl
data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  search = [
	"customRules",
	"selectedHosts"
	]
}

output "exported_configuration_text" {
  value = data.akamai_appsec_export_configuration.export.output_text
}
```

Note that you can specify a version of the configuration other than the most recent version. See the [`akamai_configuration_version`](../data-sources/appsec_configuration_version.md) data source to list the available versions. Also, you can specify other kinds of data to be exported besides `customRules` and `selectedHosts`, using any of these search fields:

* customRules
* matchTargets
* ratePolicies
* reputationProfiles
* rulesets
* securityPolicies
* selectableHosts
* selectedHosts

Save the file and run `terraform apply` to see a formatted display of the selected data.

## Adding a Hostname to the `selectedHosts` List

You can modify the list of hosts protected by a given security configuration using the [`akamai_appsec_selected_hostnames`](../data-sources/appsec_selected_hostnames.md) resource. Add the following resource block to your `akamai.tf` file, replacing `example.com` with a hostname from the list reported in the `data_akamai_appsec_export_configuration` data source example above:

```hcl
resource "akamai_appsec_selected_hostnames" "selected_hostnames_append" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  hostnames = [ "example.com" ]
  mode = "APPEND"
}

output "selected_hostnames_appended" {
  value = akamai_appsec_selected_hostnames.selected_hostnames_append.hostnames
}
```

Once you save the file and run `terraform apply`, Terraform will update the list of selected hosts and output the new list as the value `selected_hostnames_appended`. 

Note that you cannot modify a security configuration version that is currently active in staging or production, so the resource block above must specify an inactive version. Once you have completed any changes you want to make to a security configuration version, you can activate it in staging.

## Activating a Security Configuration Version

You can activate a specific version of a security configuration using the [`akamai_appsec_activations`](../resources/appsec_activations.md) resource. Add the following resource block to your `akamai.tf` file, replacing the `version` value with the number of a currently inactive version, such as the one you modified using the `akamai_appsec_selected_hostnames` resource above.

```hcl
resource "akamai_appsec_activations" "activation" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  network = "STAGING"
  notes  = "TEST Notes"
  notification_emails = [ "my_name@mycompany.com" ]
}
```

Once you save the file and run `terraform apply`, Terraform will activate the security configuration version in staging. When the activation is complete, an email will be sent to any addresses specified in the `notification_emails` list.


## Beta Features

Note that the following data sources and resources are currently in Beta, and their behavior or documentation may change in a future release:

### Data Sources
  * akamai_appsec_eval
  * akamai_appsec_eval_rule_actions
  * akamai_appsec_eval_rule_condition_exception
  * akamai_appsec_ip_geo
  * akamai_appsec_rule_actions
  * akamai_appsec_rule_condition_exception
  * akamai_appsec_penalty_box
  * akamai_appsec_security_policy_protections
  * akamai_appsec_rate_policies
  * akamai_appsec_rate_policy_actions
  * akamai_appsec_rate_protections
  * akamai_appsec_reputation_protections
  * akamai_appsec_reputation_profiles
  * akamai_appsec_reputation_profile_actions
  * akamai_appsec_rule_upgrade_details
  * akamai_appsec_slow_post
  * akamai_appsec_slowpost_protections
  * akamai_appsec_attack_group_actions
  * akamai_appsec_waf_mode
  * akamai_appsec_waf_protection
  * akamai_appsec_attack_group_condition_exception

### Resources
  * akamai_appsec_eval
  * akamai_appsec_eval_rule_action
  * akamai_appsec_eval_rule_condition_exception
  * akamai_appsec_ip_geo
  * akamai_appsec_rule_condition_exception
  * akamai_appsec_rule_action
  * akamai_appsec_penalty_box
  * akamai_appsec_security_policy_protections
  * akamai_appsec_rate_policy
  * akamai_appsec_rate_policy_action
  * akamai_appsec_rate_protection
  * akamai_appsec_reputation_protection
  * akamai_appsec_reputation_profile
  * akamai_appsec_reputation_profile_action
  * akamai_appsec_rule_upgrade
  * akamai_appsec_slow_post
  * akamai_appsec_slowpost_protection
  * akamai_appsec_attack_group_action
  * akamai_appsec_waf_mode
  * akamai_appsec_waf_protection
  * akamai_appsec_attack_group_condition_exception

