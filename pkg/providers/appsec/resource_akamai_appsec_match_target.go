package appsec

import (
	"context"
	"encoding/json"
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
func resourceMatchTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchTargetCreate,
		ReadContext:   resourceMatchTargetRead,
		UpdateContext: resourceMatchTargetUpdate,
		DeleteContext: resourceMatchTargetDelete,
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
			"match_target": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
			},
			"match_target_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceMatchTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetCreate")

	createMatchTarget := appsec.CreateMatchTargetRequest{}

	jsonpostpayload := d.Get("match_target")

	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createMatchTarget.JsonPayloadRaw = rawJSON

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createMatchTarget.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createMatchTarget.ConfigVersion = version

	postresp, err := client.CreateMatchTarget(ctx, createMatchTarget)
	if err != nil {
		logger.Errorf("calling 'createMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(postresp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("match_target", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("match_target_id", postresp.TargetID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(postresp.TargetID))

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetUpdate")

	updateMatchTarget := appsec.UpdateMatchTargetRequest{}

	jsonpostpayload := d.Get("match_target")

	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateMatchTarget.JsonPayloadRaw = rawJSON

	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateMatchTarget.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateMatchTarget.ConfigVersion = version

		targetID, errconv := strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateMatchTarget.TargetID = targetID

	} else {
		targetID, errconv := strconv.Atoi(d.Id())

		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateMatchTarget.TargetID = targetID

		jsonBody, err := json.Marshal(updateMatchTarget)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("match_target", string(jsonBody)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	resp, err := client.UpdateMatchTarget(ctx, updateMatchTarget)
	if err != nil {
		logger.Errorf("calling 'updateMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}
	jsonBody, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("match_target", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetRemove")

	removeMatchTarget := appsec.RemoveMatchTargetRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeMatchTarget.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeMatchTarget.ConfigVersion = version

		targetID, errconv := strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeMatchTarget.TargetID = targetID

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeMatchTarget.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeMatchTarget.ConfigVersion = version

		targetID, errconv := strconv.Atoi(d.Id())

		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeMatchTarget.TargetID = targetID
	}
	_, errd := client.RemoveMatchTarget(ctx, removeMatchTarget)
	if errd != nil {
		logger.Errorf("calling 'removeMatchTarget': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceMatchTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetRead")

	getMatchTarget := appsec.GetMatchTargetRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getMatchTarget.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getMatchTarget.ConfigVersion = version

		targetID, errconv := strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getMatchTarget.TargetID = targetID

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getMatchTarget.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getMatchTarget.ConfigVersion = version

		targetID, errconv := strconv.Atoi(d.Id())

		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getMatchTarget.TargetID = targetID
	}
	matchtarget, err := client.GetMatchTarget(ctx, getMatchTarget)
	if err != nil {
		logger.Errorf("calling 'getMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(matchtarget)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Warnf("calling 'getMatchTarget': JSON  %s", string(jsonBody))
	if err := d.Set("match_target", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("match_target_id", matchtarget.TargetID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("config_id", getMatchTarget.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getMatchTarget.ConfigVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%d", getMatchTarget.ConfigID, getMatchTarget.ConfigVersion, matchtarget.TargetID))

	return nil
}
