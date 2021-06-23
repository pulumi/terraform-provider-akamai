package property

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

func TestResProperty(t *testing.T) {
	// These more or less track the state of a Property in PAPI for the lifecycle tests
	type TestState struct {
		Client     *mockpapi
		Property   papi.Property
		Hostnames  []papi.Hostname
		Rules      papi.RulesUpdate
		RuleFormat string
	}

	// BehaviorFuncs can be composed to define common patterns of mock PAPI behavior (for Lifecycle tests)
	type BehaviorFunc = func(*TestState)

	// Combines many BehaviorFuncs into one
	ComposeBehaviors := func(behaviors ...BehaviorFunc) BehaviorFunc {
		return func(State *TestState) {
			for _, behave := range behaviors {
				behave(State)
			}
		}
	}

	UpdatePropertyVersionHostnames := func(PropertyID string, Version int, CnameTo string) BehaviorFunc {
		return func(State *TestState) {
			NewHostnames := []papi.Hostname{{
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            "from.test.domain",
				CnameTo:              CnameTo,
				CertProvisioningType: "DEFAULT",
			}}

			ExpectUpdatePropertyVersionHostnames(State.Client, PropertyID, "grp_0", "ctr_0", Version, NewHostnames).Once().Run(func(mock.Arguments) {
				NewResponseHostnames := []papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "from.test.domain",
					CnameTo:              CnameTo,
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_123",
					CertStatus: papi.CertStatusItem{
						ValidationCname: papi.ValidationCname{
							Hostname: "_acme-challenge.www.example.com",
							Target:   "{token}.www.example.com.akamai-domain.com",
						},
						Staging: []papi.StatusItem{{Status: "PENDING"}},
						Production: []papi.StatusItem{{
							Status: "PENDING",
						},
						},
					},
				}}
				State.Hostnames = append([]papi.Hostname{}, NewResponseHostnames...)
			})
		}
	}

	GetPropertyVersions := func(PropertyID, ContractID, GroupID string, versionItems papi.PropertyVersionItems, err error) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetPropertyVersions(State.Client, PropertyID, ContractID, GroupID, versionItems, err)
		}
	}

	GetPropertyVersionResources := func(PropertyID, GroupID, ContractID string, Version int, StagStatus, ProdStatus papi.VersionStatus) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetPropertyVersion(State.Client, PropertyID, GroupID, ContractID, Version, StagStatus, ProdStatus)
		}
	}

	GetVersionResources := func(PropertyID, ContractID, GroupID string, Version int) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetPropertyVersionHostnames(State.Client, PropertyID, GroupID, ContractID, Version, &State.Hostnames)
			ExpectGetRuleTree(State.Client, PropertyID, GroupID, ContractID, Version, &State.Rules, &State.RuleFormat)
		}
	}

	DeleteProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			ExpectRemoveProperty(State.Client, PropertyID, "ctr_0", "grp_0").Once().Run(func(mock.Arguments) {
				State.Property = papi.Property{}
				State.Rules = papi.RulesUpdate{}
				State.Hostnames = nil
				State.RuleFormat = ""
			})
		}
	}

	GetProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetProperty(State.Client, PropertyID, "grp_0", "ctr_0", &State.Property)
		}
	}

	UpdateRuleTree := func() BehaviorFunc {
		return func(State *TestState) {
			ExpectUpdateRuleTree(State.Client, "prp_0", "grp_0", "ctr_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}, "", nil)
		}
	}

	CreateProperty := func(PropertyName, PropertyID string, latestVersion int, stagingVersion, productionVersion *int) BehaviorFunc {
		return func(State *TestState) {
			ExpectCreateProperty(State.Client, PropertyName, "grp_0", "ctr_0", "prd_0", PropertyID).Run(func(mock.Arguments) {

				State.Property = papi.Property{
					PropertyName:      PropertyName,
					PropertyID:        PropertyID,
					GroupID:           "grp_0",
					ContractID:        "ctr_0",
					ProductID:         "prd_0",
					LatestVersion:     latestVersion,
					StagingVersion:    stagingVersion,
					ProductionVersion: productionVersion,
				}

				State.Rules = papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}
				State.RuleFormat = "v2020-01-01"
				GetProperty(PropertyID)(State)
				GetVersionResources(PropertyID, "ctr_0", "grp_0", 1)(State)
			}).Once()
		}
	}

	PropertyLifecycle := func(PropertyName, PropertyID, GroupID string, latestVersion, stagingVersion, productionVersion int) BehaviorFunc {
		return func(State *TestState) {
			CreateProperty(PropertyName, PropertyID, latestVersion, &stagingVersion, &productionVersion)(State)
			GetVersionResources(PropertyID, "ctr_0", "grp_0", 1)(State)
			DeleteProperty(PropertyID)(State)
		}
	}

	ImportProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			// Depending on how much of the import ID is given, the initial property lookup may not have group/contract
			ExpectGetProperty(State.Client, "prp_0", "grp_0", "", &State.Property).Maybe()
			ExpectGetProperty(State.Client, "prp_0", "", "", &State.Property).Maybe()
		}
	}

	// TestCheckFunc to verify all standard attributes (for Lifecycle tests)
	CheckAttrs := func(PropertyID, CnameTo, LatestVersion, StagingVersion, ProductionVersion, EdgeHostnameId string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_property.test", "id", PropertyID),
			resource.TestCheckResourceAttr("akamai_property.test", "hostnames.0.cname_to", CnameTo),
			resource.TestCheckResourceAttr("akamai_property.test", "hostnames.0.edge_hostname_id", EdgeHostnameId),
			resource.TestCheckResourceAttr("akamai_property.test", "latest_version", LatestVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "staging_version", StagingVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "production_version", ProductionVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "name", "test property"),
			resource.TestCheckResourceAttr("akamai_property.test", "contract_id", "ctr_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "contract", "ctr_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "group_id", "grp_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "group", "grp_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "product", "prd_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "product_id", "prd_0"),
			resource.TestCheckNoResourceAttr("akamai_property.test", "rule_warnings"),
			resource.TestCheckResourceAttr("akamai_property.test", "rules", `{"rules":{"name":"default","options":{}}}`),
		)
	}

	type StepsFunc = func(State *TestState, FixturePath string) []resource.TestStep

	// Defines standard variations of client behaviors for a Lifecycle test
	type LifecycleTestCase struct {
		Name        string
		ClientSetup BehaviorFunc
		Steps       StepsFunc
	}

	// Standard test behavior for cases where the property's latest version is deactivated in staging network
	GetLatestVersionDeactivatedInStaging := func() LifecycleTestCase {
		var stagingVersion, productionVersion *int
		stagingVersion = new(int)
		productionVersion = new(int)
		*stagingVersion = 1
		*productionVersion = 0
		LatestVersionDeactivatedInStaging := LifecycleTestCase{
			Name: "Latest version deactivated in staging",
			ClientSetup: ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0", 1, 0, 0),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusDeactivated, papi.VersionStatusInactive),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{
					{PropertyVersion: 1, StagingStatus: papi.VersionStatusActive, ProductionStatus: papi.VersionStatusInactive},
					{PropertyVersion: 2, StagingStatus: papi.VersionStatusInactive, ProductionStatus: papi.VersionStatusActive}}}, nil),
				CreateProperty("test property", "prp_0", 2, stagingVersion, productionVersion),
				UpdatePropertyVersionHostnames("prp_0", 1, "to2.test.domain"),
				UpdateRuleTree(),
				DeleteProperty("prp_0"),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
				GetVersionResources("prp_0", "ctr_0", "grp_0", 2),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusDeactivated, papi.VersionStatusInactive),
			),
			Steps: func(State *TestState, FixturePath string) []resource.TestStep {
				return []resource.TestStep{
					{
						Config:             loadFixtureString("%s/step0.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
					{
						PreConfig: func() {
							StagingVersion := 1
							State.Property.StagingVersion = &StagingVersion
						},
						Config:             loadFixtureString("%s/step1.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to2.test.domain", "2", "1", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
				}
			},
		}
		return LatestVersionDeactivatedInStaging
	}

	// Standard test behavior for cases where the property's latest version is deactivated in production network
	GetLatestVersionDeactivatedInProd := func() LifecycleTestCase {
		var stagingVersion, productionVersion *int
		stagingVersion = new(int)
		productionVersion = new(int)
		*productionVersion = 1
		LatestVersionDeactivatedInProd := LifecycleTestCase{
			Name: "Latest version is not active in production",
			ClientSetup: ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0", 1, 0, 0),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusDeactivated),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1}, {PropertyVersion: 2, ProductionStatus: papi.VersionStatusInactive}}}, nil),
				CreateProperty("test property", "prp_0", 2, stagingVersion, productionVersion),
				UpdatePropertyVersionHostnames("prp_0", 1, "to2.test.domain"),
				UpdateRuleTree(),
				DeleteProperty("prp_0"),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
				GetVersionResources("prp_0", "ctr_0", "grp_0", 2),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusDeactivated),
			),
			Steps: func(State *TestState, FixturePath string) []resource.TestStep {
				return []resource.TestStep{
					{
						Config:             loadFixtureString("%s/step0.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
					{
						PreConfig: func() {
							ProductionVersion := 1
							State.Property.ProductionVersion = &ProductionVersion
						},
						Config:             loadFixtureString("%s/step1.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to2.test.domain", "2", "0", "1", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
				}
			},
		}
		return LatestVersionDeactivatedInProd
	}

	// Standard test behavior for cases where the property's latest version is active in staging network
	GetLatestVersionActiveInStaging := func(updateruletree bool) LifecycleTestCase {
		var staging = new(int)
		*staging = 1
		LatestVersionActiveInStaging := LifecycleTestCase{
			Name: "Latest version is active in staging",
			ClientSetup: ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0", 1, 0, 0),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, StagingStatus: papi.VersionStatusInactive}, {PropertyVersion: 2, StagingStatus: papi.VersionStatusInactive}}}, nil),
				CreateProperty("test property", "prp_0", 2, staging, nil),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusActive, papi.VersionStatusInactive),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
				UpdatePropertyVersionHostnames("prp_0", 1, "to2.test.domain"),
				GetVersionResources("prp_0", "ctr_0", "grp_0", 2),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusActive, papi.VersionStatusInactive),
				DeleteProperty("prp_0"),
			),
			Steps: func(State *TestState, FixturePath string) []resource.TestStep {
				return []resource.TestStep{
					{
						Config:             loadFixtureString("%s/step0.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
					{
						PreConfig: func() {
							StagingVersion := 1
							State.Property.StagingVersion = &StagingVersion
						},
						Config:             loadFixtureString("%s/step1.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to2.test.domain", "2", "1", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
				}
			},
		}
		if updateruletree {
			LatestVersionActiveInStaging.ClientSetup = ComposeBehaviors(LatestVersionActiveInStaging.ClientSetup, UpdateRuleTree())
		}
		return LatestVersionActiveInStaging
	}

	// Standard test behavior for cases where the property's latest version is active in production network
	GetLatestVersionActiveInProd := func(updateRuleTree bool) LifecycleTestCase {
		var prodVersion = new(int)
		*prodVersion = 1
		LatestVersionActiveInProd := LifecycleTestCase{
			Name: "Latest version is active in production",
			ClientSetup: ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0", 1, 0, 0),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1}, {PropertyVersion: 2, ProductionStatus: papi.VersionStatusActive}}}, nil),
				CreateProperty("test property", "prp_0", 2, nil, prodVersion),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusActive),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
				UpdatePropertyVersionHostnames("prp_0", 1, "to2.test.domain"),
				GetVersionResources("prp_0", "ctr_0", "grp_0", 2),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusActive),
				DeleteProperty("prp_0"),
			),
			Steps: func(State *TestState, FixturePath string) []resource.TestStep {
				return []resource.TestStep{
					{
						Config:             loadFixtureString("%s/step0.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
					{
						PreConfig: func() {
							ProductionVersion := 1
							State.Property.ProductionVersion = &ProductionVersion
						},
						Config:             loadFixtureString("%s/step1.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to2.test.domain", "2", "0", "1", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
				}
			},
		}

		if updateRuleTree {
			LatestVersionActiveInProd.ClientSetup = ComposeBehaviors(LatestVersionActiveInProd.ClientSetup, UpdateRuleTree())
		}

		return LatestVersionActiveInProd
	}

	// Standard test behavior for cases where the property's latest version is not active
	GetLatestVersionNotActive := func(updateruletree bool, hostnames []string) LifecycleTestCase {
		LatestVersionNotActive := LifecycleTestCase{
			Name: "Latest version not active",
			ClientSetup: ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0", 1, 0, 0),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1}}}, nil),
				CreateProperty("test property", "prp_0", 1, nil, nil),
				DeleteProperty("prp_0"),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			),
			Steps: func(State *TestState, FixturePath string) []resource.TestStep {
				return []resource.TestStep{
					{
						Config:             loadFixtureString("%s/step0.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
					{
						Config:             loadFixtureString("%s/step1.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to2.test.domain", "1", "0", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
				}
			},
		}

		for _, host := range hostnames {
			LatestVersionNotActive.ClientSetup = ComposeBehaviors(LatestVersionNotActive.ClientSetup, UpdatePropertyVersionHostnames("prp_0", 1, host))
		}

		if updateruletree {
			LatestVersionNotActive.ClientSetup = ComposeBehaviors(LatestVersionNotActive.ClientSetup, UpdateRuleTree())
		}

		return LatestVersionNotActive
	}

	// Standard test behavior for cases where there is no diff in update
	GetNoDiff := func() LifecycleTestCase {
		var stagingVersion, productionVersion *int
		stagingVersion = new(int)
		productionVersion = new(int)
		NoDiff := LifecycleTestCase{
			Name: "No diff found in update",
			ClientSetup: ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0", 1, 0, 0),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1}}}, nil),
				CreateProperty("test property", "prp_0", 1, stagingVersion, productionVersion),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
				DeleteProperty("prp_0"),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
			),
			Steps: func(State *TestState, FixturePath string) []resource.TestStep {
				return []resource.TestStep{
					{
						Config:             loadFixtureString("%s/step0.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
					{
						Config:             loadFixtureString("%s/step1.tf", FixturePath),
						Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
						ExpectNonEmptyPlan: true,
					},
				}
			},
		}
		return NoDiff
	}

	// Run a test case to verify schema validations
	AssertConfigError := func(t *testing.T, flaw, rx string) {
		t.Helper()
		caseName := fmt.Sprintf("ConfigError/%s", flaw)

		t.Run(caseName, func(t *testing.T) {
			t.Helper()

			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/%s.tf", t.Name()),
					ExpectError: regexp.MustCompile(rx),
				}},
			})
		})
	}

	// Run a test case to verify schema attribute deprecation
	AssertDeprecated := func(t *testing.T, attribute string) {
		t.Helper()

		t.Run(fmt.Sprintf("%s attribute is deprecated", attribute), func(t *testing.T) {
			t.Helper()
			if resourceProperty().Schema[attribute].Deprecated == "" {
				t.Fatalf(`%q attribute is not marked deprecated`, attribute)
			}
		})
	}

	// Run a test case to confirm that the user is prompted to read the upgrade guide
	AssertForbiddenAttr := func(t *testing.T, fixtureName string) {
		t.Helper()

		t.Run(fmt.Sprintf("ForbiddenAttr/%s", fixtureName), func(t *testing.T) {
			t.Helper()
			client := &mockpapi{}
			client.Test(T{t})

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{{
						Config:      loadFixtureString("testdata/%s.tf", t.Name()),
						ExpectError: regexp.MustCompile("See the Akamai Terraform Upgrade Guide"),
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}

	// Run a happy-path test case that goes through a complete create-update-destroy cycle
	AssertLifecycle := func(t *testing.T, variant string, kase LifecycleTestCase) {
		t.Helper()

		fixturePrefix := fmt.Sprintf("testdata/%s/Lifecycle/%s", t.Name(), variant)
		testName := fmt.Sprintf("Lifecycle/%s/%s", variant, kase.Name)

		t.Run(testName, func(t *testing.T) {
			t.Helper()

			client := &mockpapi{}
			client.Test(T{t})
			State := &TestState{Client: client}
			kase.ClientSetup(State)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps:     kase.Steps(State, fixturePrefix),
				})
			})

			client.AssertExpectations(t)
		})
	}

	// Run a test case that verifies the resource can be imported by the given ID
	AssertImportable := func(t *testing.T, TestName, ImportID string) {
		t.Helper()

		fixturePath := fmt.Sprintf("testdata/%s/Importable/importable.tf", t.Name())
		testName := fmt.Sprintf("Importable/%s", TestName)

		t.Run(testName, func(t *testing.T) {
			t.Helper()

			client := &mockpapi{}
			client.Test(T{t})

			setup := ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0", 1, 0, 0),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusActive, StagingStatus: papi.VersionStatusActive}}}, nil),
				CreateProperty("test property", "prp_0", 1, new(int), new(int)),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusActive, papi.VersionStatusInactive),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
				DeleteProperty("prp_0"),
				ImportProperty("prp_0"),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
			)

			parameters := strings.Split(ImportID, ",")
			numberParameters := len(parameters)
			lastParameter := parameters[len(parameters)-1]
			if (numberParameters == 2 || numberParameters == 4) && !isDefaultVersion(lastParameter) {
				var ContractID, GroupID string
				if numberParameters == 4 {
					ContractID = "ctr_0"
					GroupID = "grp_0"
				}

				if numberParameters == 2 {
					setup = ComposeBehaviors(
						setup,
						GetPropertyVersionResources("prp_0", GroupID, ContractID, 1,
							papi.VersionStatusActive, papi.VersionStatusInactive),
					)
				}

				setup = ComposeBehaviors(
					setup,
					GetVersionResources("prp_0", ContractID, GroupID, 1),
					GetPropertyVersions("prp_0", ContractID, GroupID,
						papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{
							{
								PropertyVersion:  1,
								StagingStatus:    papi.VersionStatusActive,
								ProductionStatus: papi.VersionStatusActive,
							},
						}}, nil),
				)
			}

			setup(&TestState{Client: client})

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:             loadFixtureString(fixturePath),
							Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
							ExpectNonEmptyPlan: true,
						},
						{
							ImportState:       true,
							ImportStateVerify: true,
							ImportStateId:     ImportID,
							ResourceName:      "akamai_property.test",
							Config:            loadFixtureString(fixturePath),
						},
						{
							Config:             loadFixtureString(fixturePath),
							Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
							ExpectNonEmptyPlan: true,
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	}

	t.Run("invalid import ID passed", func(t *testing.T) {
		t.Helper()
		client := &mockpapi{}
		client.Test(T{t})
		ImportID := "prp_0,grp_0"
		TODO(t, "error assertion in import is impossible using provider testing framework as it only checks for errors in `apply`")
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:        loadFixtureString("testdata/TestResProperty/Importable/importable.tf"),
						ImportState:   true,
						ImportStateId: ImportID,
						ResourceName:  "akamai_property.test",
						ExpectError:   regexp.MustCompile("Either PropertyId or comma-separated list of PropertyId, contractID and groupID in that order has to be supplied in import: prp_0,grp_0"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	suppressLogging(t, func() {
		AssertConfigError(t, "name not given", `"name" is required`)
		AssertConfigError(t, "neither contract nor contract_id given", `one of .contract,contract_id. must be specified`)
		AssertConfigError(t, "both contract and contract_id given", `only one of .contract,contract_id. can be specified`)
		AssertConfigError(t, "neither group nor group_id given", `one of .group,group_id. must be specified`)
		AssertConfigError(t, "both group and group_id given", `only one of .group,group_id. can be specified`)
		AssertConfigError(t, "neither product nor product_id given", `one of product,product_id must be specified`)
		AssertConfigError(t, "both product and product_id given", `"product": conflicts with product_id`)
		AssertConfigError(t, "invalid json rules", `rules are not valid JSON`)

		AssertDeprecated(t, "contract")
		AssertDeprecated(t, "group")
		AssertDeprecated(t, "product")
		AssertDeprecated(t, "cp_code")
		AssertDeprecated(t, "contact")
		AssertDeprecated(t, "origin")
		AssertDeprecated(t, "is_secure")
		AssertDeprecated(t, "variables")

		AssertForbiddenAttr(t, "cp_code")
		AssertForbiddenAttr(t, "contact")
		AssertForbiddenAttr(t, "origin")
		AssertForbiddenAttr(t, "is_secure")
		AssertForbiddenAttr(t, "variables")

		AssertLifecycle(t, "normal", GetLatestVersionNotActive(true, []string{"to2.test.domain", "to.test.domain"}))
		AssertLifecycle(t, "normal", GetLatestVersionActiveInStaging(true))
		AssertLifecycle(t, "normal", GetLatestVersionActiveInProd(true))
		AssertLifecycle(t, "normal", GetLatestVersionDeactivatedInStaging())
		AssertLifecycle(t, "normal", GetLatestVersionDeactivatedInProd())
		AssertLifecycle(t, "contract_id without prefix", GetLatestVersionNotActive(false, []string{"to2.test.domain", "to.test.domain"}))
		AssertLifecycle(t, "contract_id without prefix", GetLatestVersionActiveInStaging(false))
		AssertLifecycle(t, "contract_id without prefix", GetLatestVersionActiveInProd(false))
		AssertLifecycle(t, "contract without prefix", GetLatestVersionNotActive(false, []string{"to2.test.domain", "to.test.domain"}))
		AssertLifecycle(t, "contract without prefix", GetLatestVersionActiveInStaging(false))
		AssertLifecycle(t, "contract without prefix", GetLatestVersionActiveInProd(false))
		AssertLifecycle(t, "group_id without prefix", GetLatestVersionNotActive(false, []string{"to2.test.domain", "to.test.domain"}))
		AssertLifecycle(t, "group_id without prefix", GetLatestVersionActiveInStaging(false))
		AssertLifecycle(t, "group_id without prefix", GetLatestVersionActiveInProd(false))
		AssertLifecycle(t, "group without prefix", GetLatestVersionNotActive(false, []string{"to2.test.domain", "to.test.domain"}))
		AssertLifecycle(t, "group without prefix", GetLatestVersionActiveInStaging(false))
		AssertLifecycle(t, "group without prefix", GetLatestVersionActiveInProd(false))
		AssertLifecycle(t, "product_id without prefix", GetLatestVersionNotActive(false, []string{"to2.test.domain", "to.test.domain"}))
		AssertLifecycle(t, "product_id without prefix", GetLatestVersionActiveInStaging(false))
		AssertLifecycle(t, "product_id without prefix", GetLatestVersionActiveInProd(false))
		AssertLifecycle(t, "product without prefix", GetLatestVersionNotActive(false, []string{"to2.test.domain", "to.test.domain"}))
		AssertLifecycle(t, "product without prefix", GetLatestVersionActiveInStaging(false))
		AssertLifecycle(t, "product without prefix", GetLatestVersionActiveInProd(false))
		AssertLifecycle(t, "no diff", GetNoDiff())
		AssertLifecycle(t, "product to product_id", GetNoDiff())
		AssertLifecycle(t, "product_id to product", GetNoDiff())

		AssertImportable(t, "property_id", "prp_0")
		AssertImportable(t, "property_id and ver_# version", "prp_0,ver_1")
		AssertImportable(t, "property_id and # version", "prp_0,1")
		AssertImportable(t, "property_id and latest", "prp_0,latest")
		AssertImportable(t, "property_id and network", "prp_0,staging")
		AssertImportable(t, "unprefixed property_id", "0")
		AssertImportable(t, "unprefixed property_id and # version", "0,1")
		AssertImportable(t, "unprefixed property_id and ver_# version", "0,ver_1")
		AssertImportable(t, "unprefixed property_id and network", "0,p")
		AssertImportable(t, "property_id and contract_id and group_id", "prp_0,ctr_0,grp_0")
		AssertImportable(t, "property_id, contract_id, group_id and empty version", "prp_0,ctr_0,grp_0,")
		AssertImportable(t, "property_id, contract_id, group_id and latest", "prp_0,ctr_0,grp_0,latest")
		AssertImportable(t, "property_id, contract_id, group_id and ver_# version", "prp_0,ctr_0,grp_0,ver_1")
		AssertImportable(t, "property_id, contract_id, group_id and # version", "prp_0,ctr_0,grp_0,1")
		AssertImportable(t, "property_id, contract_id, group_id and network", "prp_0,ctr_0,grp_0,staging")
		AssertImportable(t, "unprefixed property_id and contract_id and group_id", "0,0,0")
		AssertImportable(t, "unprefixed property_id and contract_id, group_id and # version", "0,0,0,1")
		AssertImportable(t, "unprefixed property_id and contract_id, group_id and ver_# version", "0,0,0,ver_1")
		AssertImportable(t, "unprefixed property_id and contract_id, group_id and latest", "0,0,0,latest")
		AssertImportable(t, "unprefixed property_id and contract_id, group_id and network", "0,0,0,production")

		t.Run("property is destroyed and recreated when name is changed", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			var ver1, ver2 *int
			ver1 = new(int)
			ver2 = new(int)
			*ver1 = 1
			*ver2 = 2

			setup := ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0", 1, 0, 0),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1}}}, nil),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusActive, papi.VersionStatusInactive),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
				PropertyLifecycle("renamed property", "prp_1", "grp_0", 1, 1, 1),
				GetPropertyVersionResources("prp_1", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				GetPropertyVersions("prp_1", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1}}}, nil),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1}}}, nil),
				UpdatePropertyVersionHostnames("prp_1", 1, "to2.test.domain"),
			)
			setup(&TestState{Client: client})

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:             loadFixtureString("testdata/%s-step0.tf", t.Name()),
							Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
							ExpectNonEmptyPlan: true,
						},
						{
							Config: loadFixtureString("testdata/%s-step1.tf", t.Name()),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_1"),
								resource.TestCheckResourceAttr("akamai_property.test", "name", "renamed property"),
							),
							ExpectNonEmptyPlan: true,
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error when deleting active property", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			setup := ComposeBehaviors(
				CreateProperty("test property", "prp_0", 1, new(int), new(int)),
				GetProperty("prp_0"),
				GetVersionResources("prp_0", "ctr_0", "grp_0", 1),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, "ctr_0", "grp_0"),
				UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain"),
				GetPropertyVersions("prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusActive}}}, nil),
			)
			setup(&TestState{Client: client})

			// First call to remove is not successful
			req := papi.RemovePropertyRequest{
				PropertyID: "prp_0",
				ContractID: "ctr_0",
				GroupID:    "grp_0",
			}

			err := fmt.Errorf(`cannot remove active property "prp_0"`)
			client.On("RemoveProperty", AnyCTX, req).Return(nil, err).Once()

			// Second call will be successful (TF test case requires last state to be empty or it's a failed test)
			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0").Once()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:             loadFixtureString("testdata/%s/step0.tf", t.Name()),
							Check:              CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123"),
							ExpectNonEmptyPlan: true,
						},
						{
							Config:             loadFixtureString("testdata/%s/step1.tf", t.Name()),
							ExpectError:        regexp.MustCompile(`cannot remove active property`),
							ExpectNonEmptyPlan: true,
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error when the given group is not found", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			req := papi.CreatePropertyRequest{
				ContractID: "ctr_0",
				GroupID:    "grp_0",
				Property: papi.PropertyCreate{
					ProductID:    "prd_0",
					PropertyName: "property_name",
				},
			}

			var err error = &papi.Error{
				StatusCode: 404,
				Title:      "Not Found",
				Detail:     "The system was unable to locate the requested resource",
				Type:       "https://problems.luna.akamaiapis.net/papi/v0/http/not-found",
				Instance:   "https://akaa-hqgqowhpmkw32kmt-t3owzo37wb5dkern.luna-dev.akamaiapis.net/papi/v1/properties?contractId=ctr_0\\u0026groupId=grp_0#c3fe5f9b0c4a14d1",
			}

			client.On("CreateProperty", AnyCTX, req).Return(nil, err).Once()

			// the papi GetGroups call should not return any matching group
			var Groups []*papi.Group
			ExpectGetGroups(client, &Groups).Once()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{{
						Config:      loadFixtureString("testdata/TestResProperty/Creation/property.tf"),
						ExpectError: regexp.MustCompile("group not found: grp_0"),
					}},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error when creating property with non-unique name", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			req := papi.CreatePropertyRequest{
				ContractID: "ctr_0",
				GroupID:    "grp_0",
				Property: papi.PropertyCreate{
					PropertyName: "test property",
					ProductID:    "prd_0",
				},
			}

			client.On("CreateProperty", AnyCTX, req).Return(nil, fmt.Errorf("given property name is not unique"))
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString("testdata/%s.tf", t.Name()),
							Check:       resource.TestCheckNoResourceAttr("akamai_property.test", "id"),
							ExpectError: regexp.MustCompile(`property name is not unique`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error validations when updating property with rules tree", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})
			ExpectCreateProperty(
				client, "test property", "grp_0",
				"ctr_0", "prd_0", "prp_1",
			)

			var err error = &papi.Error{
				StatusCode:   400,
				Type:         "/papi/v1/errors/validation.required_behavior",
				Title:        "Missing required behavior in default rule",
				Detail:       "In order for this property to work correctly behavior Content Provider Code needs to be present in the default section",
				Instance:     "/papi/v1/properties/prp_173136/versions/3/rules#err_100",
				BehaviorName: "cpCode",
			}
			var req = papi.UpdateRulesRequest{
				PropertyID:      "prp_1",
				ContractID:      "ctr_0",
				GroupID:         "grp_0",
				PropertyVersion: 1,
				Rules: papi.RulesUpdate{Rules: papi.Rules{
					Name: "update rule tree",
				}},
				ValidateRules: true,
			}
			client.On("UpdateRuleTree", AnyCTX, req).Return(nil, err).Once()

			ExpectRemoveProperty(client, "prp_1", "", "")
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/TestResProperty/property_update_with_validation_error_for_rules.tf"),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckNoResourceAttr("akamai_property.test", "rules")),
							ExpectError: regexp.MustCompile(`validation.required_behavior`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("validation - empty plan, when updating a property hostnames to empty", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			ExpectCreateProperty(
				client, "test property", "grp_0",
				"ctr_0", "prd_0", "prp_0",
			)

			ExpectGetPropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive)

			ExpectUpdatePropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				[]papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "terraform.provider.myu877.test.net",
					CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
					CertProvisioningType: "DEFAULT",
				}},
			).Once()

			ExpectGetPropertyVersions(client, "prp_0", "ctr_0", "grp_0", papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1}}}, nil)

			ExpectCreateProperty(client, "test property", "grp_0", "ctr_0", "prd_0", "prp_0").Run(func(mock.Arguments) {

				Property := papi.Property{
					PropertyName:  "test property",
					PropertyID:    "prp_0",
					GroupID:       "grp_0",
					ContractID:    "ctr_0",
					ProductID:     "prd_0",
					LatestVersion: 1,
				}

				Rules := papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}
				RuleFormat := "v2020-01-01"
				ExpectGetProperty(client, "prp_0", "grp_0", "ctr_0", &Property)
				ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, &[]papi.Hostname{})
				ExpectGetRuleTree(client, "prp_0", "grp_0", "ctr_0", 1, &Rules, &RuleFormat)
			}).Once()

			UpdatePropertyVersionHostnames("prp_0", 1, "to.test.domain")
			UpdateRuleTree()
			DeleteProperty("prp_0")

			ExpectGetProperty(
				client, "prp_0", "grp_0", "ctr_0",
				&papi.Property{
					PropertyID: "prp_0", GroupID: "grp_0", ContractID: "ctr_0", LatestVersion: 1,
					PropertyName: "test property",
				},
			)

			ExpectGetPropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				&[]papi.Hostname{{
					CnameFrom:            "terraform.provider.myu877.test.net",
					CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
					CertProvisioningType: "DEFAULT",
				}},
			).Times(3)

			ruleFormat := ""
			ExpectGetRuleTree(
				client, "prp_0", "grp_0", "ctr_0", 1,
				&papi.RulesUpdate{}, &ruleFormat,
			)

			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0")

			ExpectUpdatePropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				[]papi.Hostname{},
			).Once()

			ExpectGetPropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				&[]papi.Hostname{},
			).Twice()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:             loadFixtureString("testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
							Check:              resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
							ExpectNonEmptyPlan: true,
						},
						{
							Config: loadFixtureString("testdata/TestResProperty/CreationUpdateNoHostnames/update/property_update.tf"),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property.test", "hostnames.#", "0"),
							),
							ExpectError:        regexp.MustCompile("atleast one hostname required to update existing list of hostnames associated to a property"),
							ExpectNonEmptyPlan: true,
						},
					},
				})
			})
		})
	})
}
