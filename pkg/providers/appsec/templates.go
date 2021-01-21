package appsec

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/jedib0t/go-pretty/v6/table"
)

type OutputTemplates map[string]*OutputTemplate

type OutputTemplate struct {
	TemplateName   string
	TemplateType   string
	TableTitle     string
	TemplateString string
}

func GetTemplate(ots map[string]*OutputTemplate, key string) (*OutputTemplate, error) {
	if f, ok := ots[key]; ok && f != nil {
		return f, nil
	} else {
		return nil, fmt.Errorf("Error not found")
	}
}

func RenderTemplates(ots map[string]*OutputTemplate, key string, str interface{}) (string, error) {
	var ostr, tstr bytes.Buffer
	templ, ok := GetTemplate(ots, key)

	if ok == nil {

		var (
			funcs = template.FuncMap{
				"join":  strings.Join,
				"quote": func(in string) string { return fmt.Sprintf("\"%s\"", in) },
				"dash": func(in int) string {
					if in == 0 {
						return "-"
					} else {
						return strconv.Itoa(in)
					}
				},
			}
		)

		t := template.Must(template.New("").Funcs(funcs).Parse(templ.TemplateString))
		if err := t.Execute(&tstr, str); err != nil {
			return "", nil
		}

		temptype := templ.TemplateType

		if temptype == "TABULAR" {
			tbl := table.NewWriter()
			tbl.SetOutputMirror(&ostr) //os.Stdout)
			tbl.SetTitle(key)
			headers := templ.TableTitle

			headercolumns := strings.Split(headers, "|")
			trhdr := table.Row{}
			for _, header := range headercolumns {
				trhdr = append(trhdr, header)
			}
			tbl.AppendHeader(trhdr)

			ar := strings.Split(tstr.String(), ",")
			for _, recContent := range ar {
				trc := []table.Row{}
				ac := strings.Split(recContent, "|")
				tr := table.Row{}
				for _, colContent := range ac {
					tr = append(tr, colContent)
				}
				trc = append(trc, tr)
				tbl.AppendRows(trc)
			}

			tbl.Render()
		} else {
			return "\n" + tstr.String(), nil
		}

		return "\n" + ostr.String(), nil
	}
	return "", nil
}

func InitTemplates(otm map[string]*OutputTemplate) {

	otm["configuration"] = &OutputTemplate{TemplateName: "Configurations", TableTitle: "Config_id|Name|Latest_version|Version_active_in_staging|Version_active_in_production", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .Configurations}}{{if $index}},{{end}}{{.ID}}|{{.Name}}|{{.LatestVersion}}|{{.StagingVersion}}|{{.ProductionVersion}}{{end}}"}
	otm["configurationVersion"] = &OutputTemplate{TemplateName: "ConfigurationVersion", TableTitle: "Version Number|Staging Status|Production Status", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .VersionList}}{{if $index}},{{end}}{{.Version}}|{{.Staging.Status}}|{{.Production.Status}}{{end}}"}

	otm["penaltyBoxDS"] = &OutputTemplate{TemplateName: "penaltyBoxDS", TableTitle: "PenaltyBoxProtection|Action", TemplateType: "TABULAR", TemplateString: "{{.PenaltyBoxProtection}}|{{.Action}}"}

	otm["selectableHosts"] = &OutputTemplate{TemplateName: "selectableHosts", TableTitle: "Hostname", TemplateType: "TABULAR", TemplateString: "{{range .SelectableHosts}}{{.}},{{end}}"}
	otm["selectableHostsDS"] = &OutputTemplate{TemplateName: "selectableHosts", TableTitle: "Hostname|ConfigIDInProduction|ConfigNameInProduction", TemplateType: "TABULAR", TemplateString: "{{range .AvailableSet}}{{.Hostname}}|{{ dash .ConfigIDInProduction }}|{{.ConfigNameInProduction}},{{end}}"}
	otm["selectedHosts"] = &OutputTemplate{TemplateName: "selectedHosts", TableTitle: "Hostnames", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .SelectedHosts}}{{if $index}},{{end}}{{.}}{{end}}"}
	otm["selectedHostsDS"] = &OutputTemplate{TemplateName: "selectedHosts", TableTitle: "Hostnames", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .HostnameList}}{{if $index}},{{end}}{{.Hostname}}{{end}}"}
	otm["selectedHosts.tf"] = &OutputTemplate{TemplateName: "selectedHosts.tf", TableTitle: "Hostname", TemplateType: "TERRAFORM", TemplateString: "\nresource \"akamai_appsec_selected_hostnames\" \"appsecselectedhostnames\" { \n config_id = {{.ConfigID}}\n version = {{.Version}}\n hostnames = [{{  range $index, $element := .SelectedHosts }}{{if $index}},{{end}}{{quote .}}{{end}}] \n }"}
	otm["ratePolicies"] = &OutputTemplate{TemplateName: "ratePolicies", TableTitle: "ID|Policy Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .RatePolicies}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["matchTargets"] = &OutputTemplate{TemplateName: "matchTargets", TableTitle: "ID|PolicyID", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .MatchTargets.WebsiteTargets}}{{if $index}},{{end}}{{.ID}}|{{.SecurityPolicy.PolicyID}}{{end}}"}
	otm["matchTargetDS"] = &OutputTemplate{TemplateName: "matchTarget", TableTitle: "ID|PolicyID", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .MatchTargets.WebsiteTargets}}{{if $index}},{{end}}{{.TargetID}}|{{.SecurityPolicy.PolicyID}}{{end}}"}
	otm["reputationProfiles"] = &OutputTemplate{TemplateName: "reputationProfiles", TableTitle: "ID|Name(Title)", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .ReputationProfiles}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}

	otm["reputationProfilesDS"] = &OutputTemplate{TemplateName: "reputationProfilesDS", TableTitle: "ID|Name(Title)", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .ReputationProfiles}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}

	otm["reputationProfilesActions"] = &OutputTemplate{TemplateName: "reputationProfilesActions", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .ReputationProfiles}}{{if $index}},{{end}}{{.ID}}| {{.Action}}{{end}}"}

	otm["customRules"] = &OutputTemplate{TemplateName: "customRules", TableTitle: "ID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .CustomRules}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["customRuleActions"] = &OutputTemplate{TemplateName: "customRuleActions", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range .SecurityPolicies}}{{range $index, $element := .CustomRuleActions}}{{if $index}},{{end}}{{.ID}}|{{.Action}}{{end}}{{end}}"}
	otm["customRuleAction"] = &OutputTemplate{TemplateName: "customRuleAction", TableTitle: "ID|Name|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .}}{{if $index}},{{end}}{{.RuleID}}|{{.Name}} |{{.Action}}{{end}}"}
	otm["ratePolicyActions"] = &OutputTemplate{TemplateName: "ratePolicyActions", TableTitle: "ID|Ipv4Action|Ipv6Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .RatePolicyActions}}{{if $index}},{{end}}{{.ID}}| {{.Ipv4Action}}|{{.Ipv6Action}}{{end}}"}
	otm["RulesDS"] = &OutputTemplate{TemplateName: "RulesDS", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .RuleActions}}{{if $index}},{{end}}{{.ID}}|{{.Action}}{{end}}"}
	otm["EvalRulesActionsDS"] = &OutputTemplate{TemplateName: "evalRulesActions", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .RuleActions}}{{if $index}},{{end}}{{.ID}}| {{.Action}}{{end}}"}

	otm["rulesets"] = &OutputTemplate{TemplateName: "rulesets", TableTitle: "ID|Name(Title)", TemplateType: "TABULAR", TemplateString: "{{range .Rulesets}}{{range $index, $element := .Rules}}{{if $index}},{{end}}{{.ID}}| {{.Title}}{{end}}{{end}}"}
	otm["securityPolicies"] = &OutputTemplate{TemplateName: "securityPolicies", TableTitle: "ID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .SecurityPolicies}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["securityPoliciesDS"] = &OutputTemplate{TemplateName: "securityPolicies", TableTitle: "ID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .Policies}}{{if $index}},{{end}}{{.PolicyID}}|{{.PolicyName}}{{end}}"}

	otm["ruleActions"] = &OutputTemplate{TemplateName: "ruleActions", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range .SecurityPolicies}}{{range $index, $element := .WebApplicationFirewall.RuleActions}}{{if $index}},{{end}}{{.ID}}| {{.Action}}{{end}}{{end}}"}
	otm["slowPostDS"] = &OutputTemplate{TemplateName: "slowPost", TableTitle: "Action|SLOW_RATE_THRESHOLD RATE|SLOW_RATE_THRESHOLD PERIOD|DURATION_THRESHOLD TIMEOUT", TemplateType: "TABULAR", TemplateString: "{{.Action}}|{{.SlowRateThreshold.Rate}}|{{.SlowRateThreshold.Period}}|{{.DurationThreshold.Timeout}}"}
	otm["slowPost"] = &OutputTemplate{TemplateName: "slowPost", TableTitle: "Action|SLOW_RATE_THRESHOLD RATE|SLOW_RATE_THRESHOLD PERIOD|DURATION_THRESHOLD TIMEOUT", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .SecurityPolicies}}{{if $index}},{{end}}{{.SlowPost.Action}}|{{.SlowPost.DurationThreshold.Timeout}}|{{.SlowPost.SlowRateThreshold.Rate}}|{{.SlowPost.SlowRateThreshold.Period}}{{end}}"}
	otm["wafModesDS"] = &OutputTemplate{TemplateName: "wafMode", TableTitle: "Current|Mode|Eval", TemplateType: "TABULAR", TemplateString: "{{.Current}}|{{.Mode}}|{{.Eval}}"}
	otm["wafProtectionDS"] = &OutputTemplate{TemplateName: "wafProtection", TableTitle: "APIConstraints|ApplicationLayerControls|BotmanControls|NetworkLayerControls|RateControls|ReputationControls|SlowPostControls", TemplateType: "TABULAR", TemplateString: "{{.ApplyAPIConstraints}}|{{.ApplyApplicationLayerControls}}|{{.ApplyBotmanControls}}|{{.ApplyNetworkLayerControls}}|{{.ApplyRateControls}}|{{.ApplyReputationControls}}|{{.ApplySlowPostControls}}"}
	otm["AttackGroupActionDS"] = &OutputTemplate{TemplateName: "AttackGroupAction", TableTitle: "GroupID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .AttackGroupActions}}{{if $index}},{{end}}{{.Group}}| {{.Action}}{{end}}"}

	otm["rateProtectionDS"] = &OutputTemplate{TemplateName: "rateProtection", TableTitle: "APIConstraints|ApplicationLayerControls|BotmanControls|NetworkLayerControls|RateControls|ReputationControls|SlowPostControls", TemplateType: "TABULAR", TemplateString: "{{.ApplyAPIConstraints}}|{{.ApplyApplicationLayerControls}}|{{.ApplyBotmanControls}}|{{.ApplyNetworkLayerControls}}|{{.ApplyRateControls}}|{{.ApplyReputationControls}}|{{.ApplySlowPostControls}}"}
	otm["reputationProtectionDS"] = &OutputTemplate{TemplateName: "reputationProtection", TableTitle: "APIConstraints|ApplicationLayerControls|BotmanControls|NetworkLayerControls|RateControls|ReputationControls|SlowPostControls", TemplateType: "TABULAR", TemplateString: "{{.ApplyAPIConstraints}}|{{.ApplyApplicationLayerControls}}|{{.ApplyBotmanControls}}|{{.ApplyNetworkLayerControls}}|{{.ApplyRateControls}}|{{.ApplyReputationControls}}|{{.ApplySlowPostControls}}"}
	otm["slowpostProtectionDS"] = &OutputTemplate{TemplateName: "slowpostProtection", TableTitle: "APIConstraints|ApplicationLayerControls|BotmanControls|NetworkLayerControls|RateControls|ReputationControls|SlowPostControls", TemplateType: "TABULAR", TemplateString: "{{.ApplyAPIConstraints}}|{{.ApplyApplicationLayerControls}}|{{.ApplyBotmanControls}}|{{.ApplyNetworkLayerControls}}|{{.ApplyRateControls}}|{{.ApplyReputationControls}}|{{.ApplySlowPostControls}}"}
	otm["RuleConditionException"] = &OutputTemplate{TemplateName: "RuleConditionException", TableTitle: "Conditions|Exceptions", TemplateType: "TABULAR", TemplateString: "{{range $index, $element :=.Conditions}}True{{else}}False{{end}}|{{with .Exception}}True{{else}}False{{end}}"}
}
