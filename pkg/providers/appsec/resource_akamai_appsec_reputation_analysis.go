package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceReputationAnalysis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationAnalysisUpdate,
		ReadContext:   resourceReputationAnalysisRead,
		UpdateContext: resourceReputationAnalysisUpdate,
		DeleteContext: resourceReputationAnalysisDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"forward_to_http_header": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"forward_shared_ip_to_http_header_siem": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceReputationAnalysisUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisUpdate")

	updateReputationAnalysis := appsec.UpdateReputationAnalysisRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateReputationAnalysis.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateReputationAnalysis.Version = version

		policyid := s[2]
		updateReputationAnalysis.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateReputationAnalysis.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateReputationAnalysis.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateReputationAnalysis.PolicyID = policyid
	}
	forwardToHttpHeader, err := tools.GetBoolValue("forward_to_http_header", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationAnalysis.ForwardToHTTPHeader = forwardToHttpHeader

	forwardSharedIpToHttpHeaderSiem, err := tools.GetBoolValue("forward_shared_ip_to_http_header_siem", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationAnalysis.ForwardSharedIPToHTTPHeaderAndSIEM = forwardSharedIpToHttpHeaderSiem

	_, erru := client.UpdateReputationAnalysis(ctx, updateReputationAnalysis)
	if erru != nil {
		logger.Errorf("calling 'updateReputationAnalysis': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", updateReputationAnalysis.ConfigID, updateReputationAnalysis.Version, updateReputationAnalysis.PolicyID))

	return resourceReputationAnalysisRead(ctx, d, m)
}

func resourceReputationAnalysisDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisRemove")

	RemoveReputationAnalysis := appsec.RemoveReputationAnalysisRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		RemoveReputationAnalysis.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		RemoveReputationAnalysis.Version = version

		policyid := s[2]
		RemoveReputationAnalysis.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		RemoveReputationAnalysis.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		RemoveReputationAnalysis.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		RemoveReputationAnalysis.PolicyID = policyid
	}
	RemoveReputationAnalysis.ForwardToHTTPHeader = false

	RemoveReputationAnalysis.ForwardSharedIPToHTTPHeaderAndSIEM = false

	_, erru := client.RemoveReputationAnalysis(ctx, RemoveReputationAnalysis)
	if erru != nil {
		logger.Errorf("calling 'RemoveReputationAnalysis': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")

	return nil
}

func resourceReputationAnalysisRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisRead")

	getReputationAnalysis := appsec.GetReputationAnalysisRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getReputationAnalysis.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getReputationAnalysis.Version = version

		policyid := s[2]
		getReputationAnalysis.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getReputationAnalysis.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getReputationAnalysis.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getReputationAnalysis.PolicyID = policyid
	}
	_, errg := client.GetReputationAnalysis(ctx, getReputationAnalysis)
	if errg != nil {
		logger.Errorf("calling 'getReputationAnalysis': %s", errg.Error())
		return diag.FromErr(errg)
	}

	if err := d.Set("config_id", getReputationAnalysis.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getReputationAnalysis.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getReputationAnalysis.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(fmt.Sprintf("%d:%d:%s", getReputationAnalysis.ConfigID, getReputationAnalysis.Version, getReputationAnalysis.PolicyID))

	return nil
}
