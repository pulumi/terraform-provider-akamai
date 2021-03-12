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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourcePenaltyBox() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePenaltyBoxUpdate,
		ReadContext:   resourcePenaltyBoxRead,
		UpdateContext: resourcePenaltyBoxUpdate,
		DeleteContext: resourcePenaltyBoxDelete,
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
			"penalty_box_protection": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"penalty_box_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Deny,
					Alert,
					None,
				}, false),
			},
		},
	}
}

func resourcePenaltyBoxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxRead")

	getPenaltyBox := appsec.GetPenaltyBoxRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getPenaltyBox.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getPenaltyBox.Version = version

		policyid := s[2]
		getPenaltyBox.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getPenaltyBox.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getPenaltyBox.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getPenaltyBox.PolicyID = policyid
	}
	penaltybox, err := client.GetPenaltyBox(ctx, getPenaltyBox)
	if err != nil {
		logger.Errorf("calling 'getPenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getPenaltyBox.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getPenaltyBox.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getPenaltyBox.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("penalty_box_protection", penaltybox.PenaltyBoxProtection); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("penalty_box_action", penaltybox.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", getPenaltyBox.ConfigID, getPenaltyBox.Version, getPenaltyBox.PolicyID))

	return nil
}

func resourcePenaltyBoxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxRemove")

	removePenaltyBox := appsec.UpdatePenaltyBoxRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removePenaltyBox.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removePenaltyBox.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removePenaltyBox.PolicyID = policyid

	removePenaltyBox.Action = "none"

	removePenaltyBox.PenaltyBoxProtection = false

	_, errd := client.UpdatePenaltyBox(ctx, removePenaltyBox)
	if errd != nil {
		logger.Errorf("calling 'removePenaltyBox': %s", errd.Error())
		return diag.FromErr(errd)
	}
	d.SetId("")

	return nil
}

func resourcePenaltyBoxUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxUpdate")

	updatePenaltyBox := appsec.UpdatePenaltyBoxRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePenaltyBox.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePenaltyBox.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePenaltyBox.PolicyID = policyid

	penaltyboxaction, err := tools.GetStringValue("penalty_box_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePenaltyBox.Action = penaltyboxaction

	penaltyboxprotection, err := tools.GetBoolValue("penalty_box_protection", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePenaltyBox.PenaltyBoxProtection = penaltyboxprotection

	_, erru := client.UpdatePenaltyBox(ctx, updatePenaltyBox)
	if erru != nil {
		logger.Errorf("calling 'updatePenaltyBox': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourcePenaltyBoxRead(ctx, d, m)
}
