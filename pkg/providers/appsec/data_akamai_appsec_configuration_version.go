package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigurationVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationVersionRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"latest_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"staging_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"production_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceConfigurationVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceConfigurationVersionRead")

	getConfigurationVersion := appsec.GetConfigurationVersionsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getConfigurationVersion.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getConfigurationVersion.ConfigVersion = version

	configurationversion, err := client.GetConfigurationVersions(ctx, getConfigurationVersion)
	if err != nil {
		logger.Errorf("calling 'getConfigurationVersion': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("latest_version", configurationversion.LastCreatedVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	for _, configval := range configurationversion.VersionList {

		if configval.Version == version {

			if err := d.Set("config_id", configval.ConfigID); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}

			if err := d.Set("staging_status", configval.Staging.Status); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}

			if err := d.Set("production_status", configval.Production.Status); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}

			d.SetId(strconv.Itoa(configval.ConfigID))
		}
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "configurationVersion", configurationversion)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	d.SetId(strconv.Itoa(configurationversion.ConfigID))

	return nil
}
