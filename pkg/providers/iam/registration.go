// +build all iam

package iam

import (
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/registry"
)

func init() {
	registry.RegisterProvider(&provider{})
}
