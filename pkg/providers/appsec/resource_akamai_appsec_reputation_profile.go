package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceReputationProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationProfileCreate,
		ReadContext:   resourceReputationProfileRead,
		UpdateContext: resourceReputationProfileUpdate,
		DeleteContext: resourceReputationProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"reputation_profile": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentReputationProfileDiffs,
			},
			"reputation_profile_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceReputationProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileCreate")
	logger.Debug("!!! in resourceReputationProfileCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProfile", m)
	jsonpostpayload, err := tools.GetStringValue("reputation_profile", d)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonPayloadRaw := []byte(jsonpostpayload)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createReputationProfile := appsec.CreateReputationProfileRequest{}
	createReputationProfile.ConfigID = configid
	createReputationProfile.ConfigVersion = version
	createReputationProfile.JsonPayloadRaw = rawJSON

	response, err := client.CreateReputationProfile(ctx, createReputationProfile)
	if err != nil {
		logger.Errorf("calling 'CreateReputationProfile': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", createReputationProfile.ConfigID, response.ID))

	return resourceReputationProfileRead(ctx, d, m)
}

func resourceReputationProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileRead")
	logger.Debug("!!! in resourceReputationProfileRead")

	idParts, err := splitID(d.Id(), 2, "configid:reputationprofileid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	reputationProfileID, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	reputationProfileRequest := appsec.GetReputationProfileRequest{}
	reputationProfileRequest.ConfigID = configid
	reputationProfileRequest.ConfigVersion = getLatestConfigVersion(ctx, configid, m)
	reputationProfileRequest.ReputationProfileId = reputationProfileID

	reputationProfileResponse, err := client.GetReputationProfile(ctx, reputationProfileRequest)
	if err != nil {
		logger.Errorf("calling 'getReputationProfile': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("reputation_profile_id", reputationProfileID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	jsonBody, err := json.Marshal(reputationProfileResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("reputation_profile", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceReputationProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileUpdate")
	logger.Debug("!!! in resourceReputationProfileUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:reputationprofileid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	reputationProfileID, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}
	jsonpostpayload, err := tools.GetStringValue("reputation_profile", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProfile", m)
	jsonPayloadRaw := []byte(jsonpostpayload)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateReputationProfile := appsec.UpdateReputationProfileRequest{}
	updateReputationProfile.ConfigID = configid
	updateReputationProfile.ConfigVersion = version
	updateReputationProfile.ReputationProfileId = reputationProfileID
	updateReputationProfile.JsonPayloadRaw = rawJSON

	_, err = client.UpdateReputationProfile(ctx, updateReputationProfile)
	if err != nil {
		logger.Errorf("calling 'updateReputationProfile': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceReputationProfileRead(ctx, d, m)
}

func resourceReputationProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileDelete")
	logger.Debug("!!! in resourceReputationProfileDelete")

	idParts, err := splitID(d.Id(), 2, "configid:reputationprofileid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProfile", m)
	reputationProfileID, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	deleteReputationProfile := appsec.RemoveReputationProfileRequest{}
	deleteReputationProfile.ConfigID = configid
	deleteReputationProfile.ConfigVersion = version
	deleteReputationProfile.ReputationProfileId = reputationProfileID

	_, errd := client.RemoveReputationProfile(ctx, deleteReputationProfile)
	if errd != nil {
		logger.Errorf("calling 'removeReputationProfile': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}
