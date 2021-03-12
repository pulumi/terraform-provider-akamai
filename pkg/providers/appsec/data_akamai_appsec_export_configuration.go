package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExportConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExportConfigurationRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"search": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON Export representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceExportConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceExportConfigurationRead")

	getExportConfiguration := appsec.GetExportConfigurationsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getExportConfiguration.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getExportConfiguration.Version = version

	exportconfiguration, err := client.GetExportConfigurations(ctx, getExportConfiguration)
	if err != nil {
		logger.Errorf("calling 'getExportConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(exportconfiguration)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	searchlist, ok := d.GetOk("search")
	if ok {
		ots := OutputTemplates{}
		InitTemplates(ots)

		var outputtextresult string

		for _, h := range searchlist.([]interface{}) {
			outputtext, err := RenderTemplates(ots, h.(string), exportconfiguration)
			if err == nil {
				outputtextresult = outputtextresult + outputtext
			}
		}

		if len(outputtextresult) > 0 {
			if err := d.Set("output_text", outputtextresult); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}
		}
	}
	d.SetId(strconv.Itoa(exportconfiguration.ConfigID))

	return nil
}
