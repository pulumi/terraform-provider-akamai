package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) resUser() *schema.Resource {
	validateAuthGrantJS := func(v interface{}, _ cty.Path) diag.Diagnostics {
		js := []byte(v.(string))
		if len(js) == 0 {
			return nil
		}

		var AuthGrants []iam.AuthGrant
		if err := json.Unmarshal(js, &AuthGrants); err != nil {
			return diag.Errorf("auth_grants_json is not valid: %s", err)
		}

		if len(AuthGrants) == 0 {
			return diag.Errorf("auth_grants_json must contain at least one entry")
		}

		return nil
	}

	stateAuthGrantsJS := func(v interface{}) string {
		js := []byte(v.(string))
		if len(js) == 0 {
			return ""
		}

		var AuthGrants []iam.AuthGrant
		if err := json.Unmarshal(js, &AuthGrants); err != nil {
			panic(fmt.Sprintf(`"auth_grants": %q is not valid: %s`, v.(string), err))
		}

		var AuthGrantsJSON []byte
		AuthGrantsJSON, err := json.Marshal(AuthGrants)
		if err != nil {
			panic(fmt.Sprintf(`"auth_grants": %q is not valid: %s`, v.(string), err))
		}

		return string(AuthGrantsJSON)
	}

	suppressAuthGrantsJS := func(k, old, new string, d *schema.ResourceData) bool {
		var Old []iam.AuthGrant
		if len(old) > 0 {
			if err := json.Unmarshal([]byte(old), &Old); err != nil {
				panic(fmt.Sprintf("previous value for %q: %q is not valid: %s", k, old, err))
			}
		}

		var New []iam.AuthGrant
		if len(new) > 0 {
			if err := json.Unmarshal([]byte(new), &New); err != nil {
				panic(fmt.Sprintf("new value for %q: %q is not valid: %s", k, new, err))
			}
		}

		return cmp.Equal(Old, New, cmpopts.EquateEmpty())
	}

	validatePhone := func(v interface{}, _ cty.Path) diag.Diagnostics {
		ph := v.(string)
		if v.(string) == "" {
			return nil
		}

		if !regexp.MustCompile(`^\(\d{3}\) \d{3}-\d{4}$`).MatchString(ph) {
			return diag.Errorf(`Phone number must be in the form: (###) ###-####`)
		}

		return nil
	}

	statePhone := func(v interface{}) string {
		return canonicalPhone(v.(string))
	}

	suppressPhone := func(k, old, new string, d *schema.ResourceData) bool {
		old = regexp.MustCompile(`[^0-9]+`).ReplaceAllLiteralString(old, "")
		new = regexp.MustCompile(`[^0-9]+`).ReplaceAllLiteralString(new, "")
		return old == new
	}

	suppressEmail := func(k, old, new string, d *schema.ResourceData) bool {
		return strings.ToLower(old) == strings.ToLower(new)
	}

	stateEmail := func(v interface{}) string {
		return strings.ToLower(v.(string))
	}

	return &schema.Resource{
		Description:   "Manage a user in your account",
		CreateContext: p.tfCRUD("res:User:Create", p.resUserCreate),
		ReadContext:   p.tfCRUD("res:User:Read", p.resUserRead),
		UpdateContext: p.tfCRUD("res:User:Update", p.resUserUpdate),
		DeleteContext: p.tfCRUD("res:User:Delete", p.resUserDelete),
		Importer:      p.tfImporter("res:User:Import", schema.ImportStatePassthroughContext),
		Schema: map[string]*schema.Schema{
			// Inputs - Required
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's first name",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's surname",
			},
			"email": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The user's email address",
				StateFunc:        stateEmail,
				DiffSuppressFunc: suppressEmail,
			},
			"country": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "As part of the user's location, the value can be any that are available from the view-supported-countries operation",
			},
			"phone": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The user's main phone number",
				ValidateDiagFunc: validatePhone,
				DiffSuppressFunc: suppressPhone,
				StateFunc:        statePhone,
			},
			"enable_tfa": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates whether two-factor authentication is allowed",
			},
			"auth_grants_json": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "A user's per-group role assignments, in JSON form",
				ValidateDiagFunc: validateAuthGrantJS,
				DiffSuppressFunc: suppressAuthGrantsJS,
				StateFunc:        stateAuthGrantsJS,
			},

			// Inputs - Optional
			"contact_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To help characterize the user, the value can be any that are available from the view-contact-types operation",
			},
			"job_title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's position at your company",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user's time zone. The value can be any that are available from the view-time-zones operation",
			},
			"secondary_email": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The user's secondary email address",
				StateFunc:        stateEmail,
				DiffSuppressFunc: suppressEmail,
			},
			"mobile_phone": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The user's mobile phone number",
				ValidateDiagFunc: validatePhone,
				DiffSuppressFunc: suppressPhone,
				StateFunc:        statePhone,
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user's street address",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's city",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's state",
			},
			"zip_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's five-digit ZIP code",
			},
			"preferred_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The value can be any that are available from the view-languages operation",
			},
			"session_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of seconds it takes for the user's Control Center session to time out if there hasn't been any activity",
			},

			// Purely computed
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A user's `loginId`. Typically, a user's email address",
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The user's lock status",
			},
			"last_login": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the user last logged in",
			},
			"password_expired_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a user's password expires",
			},
			"tfa_configured": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether two-factor authentication is configured",
			},
			"email_update_pending": {
				Type:     schema.TypeBool,
				Computed: true,
				// Description: "TODO", // 🤷‍♂️ Couldn't find this in docs or service descriptors
			},
		},
	}
}

func (p *provider) resUserCreate(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	AuthGrantsJSON := []byte(d.Get("auth_grants_json").(string))

	var AuthGrants []iam.AuthGrant
	if len(AuthGrantsJSON) > 0 {
		if err := json.Unmarshal(AuthGrantsJSON, &AuthGrants); err != nil {
			logger.WithError(err).Errorf("auth_grants is not valid")
			return diag.Errorf("auth_grants is not valid: %s", err)
		}
	}

	BasicUser := iam.UserBasicInfo{
		FirstName:         d.Get("first_name").(string),
		LastName:          d.Get("last_name").(string),
		UserName:          d.Get("user_name").(string),
		Email:             d.Get("email").(string),
		Phone:             d.Get("phone").(string),
		TimeZone:          d.Get("time_zone").(string),
		JobTitle:          d.Get("job_title").(string),
		TFAEnabled:        d.Get("enable_tfa").(bool),
		SecondaryEmail:    d.Get("secondary_email").(string),
		MobilePhone:       d.Get("mobile_phone").(string),
		Address:           d.Get("address").(string),
		City:              d.Get("city").(string),
		State:             d.Get("state").(string),
		ZipCode:           d.Get("zip_code").(string),
		Country:           d.Get("country").(string),
		ContactType:       d.Get("contact_type").(string),
		PreferredLanguage: d.Get("preferred_language").(string),
	}

	if st, ok := d.GetOk("session_timeout"); ok {
		SessionTimeout := st.(int)
		BasicUser.SessionTimeOut = &SessionTimeout
	}

	User, err := p.client.CreateUser(ctx, iam.CreateUserRequest{
		User:       BasicUser,
		AuthGrants: AuthGrants,
		SendEmail:  true,
		Notifications: iam.UserNotifications{
			Options: iam.UserNotificationOptions{
				Proactive: []string{},
				Upgrade:   []string{},
			},
		},
	})
	if err != nil {
		logger.WithError(err).Errorf("failed to create user")
		return diag.Errorf("failed to create user: %s\n%s", err, resUserErrorAdvice(err))
	}

	d.SetId(User.IdentityID)
	return p.resUserRead(ctx, d, nil)
}

func (p *provider) resUserRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	req := iam.GetUserRequest{
		IdentityID: d.Id(),
		AuthGrants: true,
	}

	User, err := p.client.GetUser(ctx, req)
	if err != nil {
		logger.WithError(err).Errorf("failed to fetch user")
		return diag.Errorf("failed to fetch user: %s", err)
	}

	if User.SessionTimeOut == nil {
		SessionTimeOut := 0
		User.SessionTimeOut = &SessionTimeOut
	}

	var AuthGrantsJSON []byte
	if len(User.AuthGrants) > 0 {
		AuthGrantsJSON, err = json.Marshal(User.AuthGrants)
		if err != nil {
			logger.WithError(err).Error("could not marshal AuthGrants")
			return diag.Errorf("could not marshal AuthGrants: %s", err)
		}
	}

	err = tools.SetAttrs(d, map[string]interface{}{
		"first_name":             User.FirstName,
		"last_name":              User.LastName,
		"user_name":              User.UserName,
		"email":                  User.Email,
		"phone":                  canonicalPhone(User.Phone),
		"time_zone":              User.TimeZone,
		"job_title":              User.JobTitle,
		"enable_tfa":             User.TFAEnabled,
		"secondary_email":        User.SecondaryEmail,
		"mobile_phone":           canonicalPhone(User.MobilePhone),
		"address":                User.Address,
		"city":                   User.City,
		"state":                  User.State,
		"zip_code":               User.ZipCode,
		"country":                User.Country,
		"contact_type":           User.ContactType,
		"preferred_language":     User.PreferredLanguage,
		"is_locked":              User.IsLocked,
		"last_login":             User.LastLoginDate,
		"password_expired_after": User.PasswordExpiryDate,
		"tfa_configured":         User.TFAConfigured,
		"email_update_pending":   User.EmailUpdatePending,
		"session_timeout":        *User.SessionTimeOut,
		"auth_grants_json":       string(AuthGrantsJSON),
	})
	if err != nil {
		logger.WithError(err).Error("could not save attributes to state")
		return diag.Errorf("could not save attributes to state: %s", err)
	}

	return nil
}

func (p *provider) resUserUpdate(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	if d.HasChange("email") {
		d.Partial(true)
		err := fmt.Errorf("cannot change email address")
		logger.WithError(err).Errorf("failed to update user")
		return diag.Errorf("failed to update user: %s", err)
	}

	var needRead bool

	// Basic Info
	updateBasicInfo := d.HasChanges(
		"first_name",
		"last_name",
		"phone",
		"time_zone",
		"job_title",
		"enable_tfa",
		"secondary_email",
		"mobile_phone",
		"address",
		"city",
		"state",
		"zip_code",
		"country",
		"contact_type",
		"preferred_language",
		"session_timeout",
	)
	if updateBasicInfo {
		BasicUser := iam.UserBasicInfo{
			FirstName:         d.Get("first_name").(string),
			LastName:          d.Get("last_name").(string),
			UserName:          d.Get("user_name").(string),
			Email:             d.Get("email").(string),
			Phone:             d.Get("phone").(string),
			TimeZone:          d.Get("time_zone").(string),
			JobTitle:          d.Get("job_title").(string),
			TFAEnabled:        d.Get("enable_tfa").(bool),
			SecondaryEmail:    d.Get("secondary_email").(string),
			MobilePhone:       d.Get("mobile_phone").(string),
			Address:           d.Get("address").(string),
			City:              d.Get("city").(string),
			State:             d.Get("state").(string),
			ZipCode:           d.Get("zip_code").(string),
			Country:           d.Get("country").(string),
			ContactType:       d.Get("contact_type").(string),
			PreferredLanguage: d.Get("preferred_language").(string),
		}

		if st, ok := d.GetOk("session_timeout"); ok {
			SessionTimeout := st.(int)
			BasicUser.SessionTimeOut = &SessionTimeout
		}

		req := iam.UpdateUserInfoRequest{
			IdentityID: d.Id(),
			User:       BasicUser,
		}
		if _, err := p.client.UpdateUserInfo(ctx, req); err != nil {
			d.Partial(true)
			logger.WithError(err).Errorf("failed to update user")
			return diag.Errorf("failed to update user: %s\n%s", err, resUserErrorAdvice(err))
		}

		needRead = true
	}

	// AuthGrants
	if d.HasChange("auth_grants_json") {
		var AuthGrants []iam.AuthGrant

		AuthGrantsJSON := []byte(d.Get("auth_grants_json").(string))
		if len(AuthGrantsJSON) > 0 {
			if err := json.Unmarshal(AuthGrantsJSON, &AuthGrants); err != nil {
				d.Partial(true)
				logger.WithError(err).Errorf("auth_grants is not valid")
				return diag.Errorf("auth_grants is not valid: %s", err)
			}
		}

		req := iam.UpdateUserAuthGrantsRequest{
			IdentityID: d.Id(),
			AuthGrants: AuthGrants,
		}
		if _, err := p.client.UpdateUserAuthGrants(ctx, req); err != nil {
			d.Partial(true)
			logger.WithError(err).Errorf("failed to update user AuthGrants")
			return diag.Errorf("failed to update user AuthGrants: %s", err)
		}

		needRead = true
	}

	if needRead {
		d.Partial(false)
		return p.resUserRead(ctx, d, nil)
	}

	d.Partial(false)
	return nil
}

func (p *provider) resUserDelete(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	if err := p.client.RemoveUser(ctx, iam.RemoveUserRequest{IdentityID: d.Id()}); err != nil {
		logger.WithError(err).Error("could not remove user")
		return diag.Errorf("could not remove user: %s", err)
	}

	return nil
}

func resUserErrorAdvice(e error) string {
	switch {
	case regexp.MustCompile(`\b(preferredLanguage|[pP]referred [lL]anguage)\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_supported_langs" data source to get possible values for "preferred_language"`

	case regexp.MustCompile(`\b(contactType|[cC]ontact [tT]ype)\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_contact_types" data source to get possible values for "contact_type"`

	case regexp.MustCompile(`\b[cC]ountry\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_countries" data source to get possible values for "country"`

	case regexp.MustCompile(`\b(sessionTimeOut|[sS]ession [tT]ime ?[oO]ut)\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_timeout_policies" data source to get possible values for "session_timeout"`

	case regexp.MustCompile(`\b[sS]tate\b`).FindStringIndex(e.Error()) != nil:
		return `Tip: Use the "akamai_iam_states" data source to get possible values for "state"`
	}

	return ""
}

func canonicalPhone(in string) string {
	ph := regexp.MustCompile(`[^0-9]+`).ReplaceAllLiteralString(in, "")
	if len(ph) < 10 {
		return in
	}

	return fmt.Sprintf("(%s) %s-%s", ph[0:3], ph[3:6], ph[6:10])
}
