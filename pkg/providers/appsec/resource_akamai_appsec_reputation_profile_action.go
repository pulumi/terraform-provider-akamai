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
func resourceReputationProfileAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationProfileActionUpdate,
		ReadContext:   resourceReputationProfileActionRead,
		UpdateContext: resourceReputationProfileActionUpdate,
		DeleteContext: resourceReputationProfileActionDelete,
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
			"reputation_profile_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
		},
	}
}

func resourceReputationProfileActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionRead")

	getReputationProfileAction := appsec.GetReputationProfileActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getReputationProfileAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getReputationProfileAction.Version = version

		policyid := s[2]
		getReputationProfileAction.PolicyID = policyid

		reputationprofileid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getReputationProfileAction.ReputationProfileID = reputationprofileid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getReputationProfileAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getReputationProfileAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getReputationProfileAction.PolicyID = policyid

		reputationprofileid, err := tools.GetIntValue("reputation_profile_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getReputationProfileAction.ReputationProfileID = reputationprofileid
	}
	resp, errr := client.GetReputationProfileAction(ctx, getReputationProfileAction)
	if errr != nil {
		logger.Errorf("calling 'getReputationProfileAction': %s", errr.Error())
		return diag.FromErr(errr)
	}

	if err := d.Set("action", resp.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", getReputationProfileAction.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getReputationProfileAction.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getReputationProfileAction.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s:%d", getReputationProfileAction.ConfigID, getReputationProfileAction.Version, getReputationProfileAction.PolicyID, getReputationProfileAction.ReputationProfileID))

	return nil
}

func resourceReputationProfileActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionRemove")

	removeReputationProfileAction := appsec.UpdateReputationProfileActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeReputationProfileAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeReputationProfileAction.Version = version

		policyid := s[2]
		removeReputationProfileAction.PolicyID = policyid

		reputationprofileid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeReputationProfileAction.ReputationProfileID = reputationprofileid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeReputationProfileAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeReputationProfileAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeReputationProfileAction.PolicyID = policyid

		reputationprofileid, err := tools.GetIntValue("reputation_profile_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeReputationProfileAction.ReputationProfileID = reputationprofileid
	}
	removeReputationProfileAction.Action = "none"

	_, errd := client.UpdateReputationProfileAction(ctx, removeReputationProfileAction)
	if errd != nil {
		logger.Errorf("calling 'removeReputationProfileAction': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceReputationProfileActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionUpdate")

	updateReputationProfileAction := appsec.UpdateReputationProfileActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateReputationProfileAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateReputationProfileAction.Version = version

		policyid := s[2]
		updateReputationProfileAction.PolicyID = policyid

		reputationprofileid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateReputationProfileAction.ReputationProfileID = reputationprofileid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateReputationProfileAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateReputationProfileAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateReputationProfileAction.PolicyID = policyid

		reputationprofileid, err := tools.GetIntValue("reputation_profile_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateReputationProfileAction.ReputationProfileID = reputationprofileid
	}
	action, err := tools.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationProfileAction.Action = action

	_, erru := client.UpdateReputationProfileAction(ctx, updateReputationProfileAction)
	if erru != nil {
		logger.Errorf("calling 'updateReputationProfileAction': %s", erru.Error())
		return diag.FromErr(erru)
	}
	d.SetId(fmt.Sprintf("%d:%d:%s:%d", updateReputationProfileAction.ConfigID, updateReputationProfileAction.Version, updateReputationProfileAction.PolicyID, updateReputationProfileAction.ReputationProfileID))

	return resourceReputationProfileActionRead(ctx, d, m)
}
