package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiVersionNotes_data_basic(t *testing.T) {
	t.Run("match by VersionNotes ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetVersionNotesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSVersionNotes/VersionNotes.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetVersionNotes",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetVersionNotesRequest{ConfigID: 43253, Version: 7},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSVersionNotes/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_version_notes.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
