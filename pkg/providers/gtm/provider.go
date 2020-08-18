package gtm

import (
	"context"
	"log"
	"sync"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider
	}
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider() akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"gtm_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("gtm_section"),
			},
			"gtm": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("gtm"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_gtm_default_datacenter": dataSourceGTMDefaultDatacenter(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_gtm_domain":     resourceGTMv1Domain(),
			"akamai_gtm_datacenter": resourceGTMv1Datacenter(),
			"akamai_gtm_property":   resourceGTMv1Property(),
			"akamai_gtm_resource":   resourceGTMv1Resource(),
			"akamai_gtm_cidrmap":    resourceGTMv1Cidrmap(),
			"akamai_gtm_geomap":     resourceGTMv1Geomap(),
			"akamai_gtm_asmap":      resourceGTMv1ASmap(),
		},
	}
	//ConfigureFunc: providerConfigure,
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}
	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	log.Printf("[DEBUG] START providerConfigure  %s\n", terraformVersion)
	cfg, err := getConfigGTMV1Service(d)
	if err != nil {
		return nil, nil
	}

	return cfg, nil
}

type resourceData interface {
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
}

type set interface {
	List() []interface{}
}

func getConfigGTMV1Service(d resourceData) (*edgegrid.Config, error) {
	var GTMv1Config edgegrid.Config
	var err error
	if _, ok := d.GetOk("gtm"); ok {
		config := d.Get("gtm").(set).List()[0].(map[string]interface{})

		GTMv1Config = edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
		}

		gtm.Init(GTMv1Config)
		edgegrid.SetupLogging()
		return &GTMv1Config, nil
	}

	edgerc := d.Get("edgerc").(string)
	section := d.Get("gtm_section").(string)
	if section == "" {
		section = d.Get("section").(string)
	}
	GTMv1Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	gtm.Init(GTMv1Config)
	return &GTMv1Config, nil
}

func (p *provider) Name() string {
	return "gtm"
}

func (p *provider) Version() string {
	return "v0.8.3"
}

func (p *provider) Schema() map[string]*schema.Schema {
	return p.Provider.Schema
}

func (p *provider) Resources() map[string]*schema.Resource {
	return p.Provider.ResourcesMap
}

func (p *provider) DataSources() map[string]*schema.Resource {
	return p.Provider.DataSourcesMap
}

func (p *provider) Configure(ctx context.Context, log hclog.Logger, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	state, err := p.ConfigureFunc(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return state, nil
}
