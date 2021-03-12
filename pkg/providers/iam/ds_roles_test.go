package iam

import (
	"errors"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
)

func TestDSRoles(t *testing.T) {
	t.Parallel()

	t.Run("happy path/no args", func(t *testing.T) {
		t.Parallel()

		roles := []iam.Role{{
			RoleName:        "test role name",
			RoleID:          100,
			RoleDescription: "role description",
			RoleType:        iam.RoleTypeStandard,
			CreatedBy:       "creator@akamai.net",
			CreatedDate:     "2020-01-01T00:00:00Z",
			ModifiedBy:      "modifier@akamai.net",
			ModifiedDate:    "2020-01-01T00:00:00Z",
		}}

		req := iam.ListRolesRequest{}

		client := &IAM{}
		client.Test(test.TattleT{T: t})
		client.On("ListRoles", mock.Anything, req).Return(roles, nil)

		p := provider{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: test.Fixture("testdata/%s.tf", t.Name()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.akamai_iam_roles.test", "id"),
						resource.TestCheckResourceAttr("data.akamai_iam_roles.test", "roles.0.name", "test role name"),
						resource.TestCheckResourceAttr("data.akamai_iam_roles.test", "roles.0.role_id", "100"),
						resource.TestCheckResourceAttr("data.akamai_iam_roles.test", "roles.0.description", "role description"),
						resource.TestCheckResourceAttr("data.akamai_iam_roles.test", "roles.0.type", string(iam.RoleTypeStandard)),
						resource.TestCheckResourceAttr("data.akamai_iam_roles.test", "roles.0.time_created", "2020-01-01T00:00:00Z"),
						resource.TestCheckResourceAttr("data.akamai_iam_roles.test", "roles.0.time_modified", "2020-01-01T00:00:00Z"),
						resource.TestCheckResourceAttr("data.akamai_iam_roles.test", "roles.0.created_by", "creator@akamai.net"),
						resource.TestCheckResourceAttr("data.akamai_iam_roles.test", "roles.0.modified_by", "modifier@akamai.net"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})

	t.Run("fail path", func(t *testing.T) {
		t.Parallel()

		req := iam.ListRolesRequest{}

		client := &IAM{}
		client.Test(test.TattleT{T: t})
		client.On("ListRoles", mock.Anything, req).Return(nil, errors.New("failed to get roles"))

		p := provider{}
		p.SetCache(metaCache{})
		p.SetIAM(client)

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: p.ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config:      test.Fixture("testdata/%s/step0.tf", t.Name()),
					ExpectError: regexp.MustCompile(`failed to get roles`),
				},
			},
		})

		client.AssertExpectations(t)
	})
}
