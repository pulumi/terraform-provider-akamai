package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsCountries() *schema.Resource {
	return &schema.Resource{
		Description: "List all the possible countries that Akamai supports",
		ReadContext: p.tfCRUD("ds:Countries:Read", p.dsCountriesRead),
		Schema: map[string]*schema.Schema{
			"countries": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Supported countries",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func (p *provider) dsCountriesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	logger.Debug("Fetching supported countries")
	res, err := p.client.SupportedCountries(ctx)
	if err != nil {
		logger.WithError(err).Error("Could not get supported countries")
		return diag.FromErr(err)
	}

	countries := []interface{}{}
	for _, country := range res {
		countries = append(countries, country)
	}

	if err := d.Set("countries", countries); err != nil {
		logger.WithError(err).Error("Could not set countries in state")
		return diag.FromErr(err)
	}

	d.SetId("countries")
	return nil
}
