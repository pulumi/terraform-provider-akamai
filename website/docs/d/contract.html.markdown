---
layout: "akamai"
page_title: "Akamai: contract"
sidebar_current: "docs-akamai-data-contract"
description: |-
 Contract
---

# akamai_contract

Use `akamai_contract` data source to retrieve a group id.

## Example Usage

### Basic usage:

```hcl
data "akamai_contract" "example" {
     group = "group name"
}

resource "akamai_property" "example" {
    contract = "${data.akamai_contract.example.id}"
    ...
}
```

## Argument Reference

The following arguments are supported:

* `group` — (Optional) The group within which the contract can be found.

## Attributes Reference

The following are the return attributes:

* `id` — The contract ID.
