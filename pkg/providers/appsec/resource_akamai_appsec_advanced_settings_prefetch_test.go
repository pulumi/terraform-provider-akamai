package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAdvancedSettingsPrefetch_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPrefetch ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateAdvancedSettingsPrefetchResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAdvancedSettingsPrefetchResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetAdvancedSettingsPrefetch",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAdvancedSettingsPrefetchRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateAdvancedSettingsPrefetch",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAdvancedSettingsPrefetchRequest{ConfigID: 43253, Version: 7, AllExtensions: true, EnableAppLayer: false, EnableRateControls: false, Extensions: []string{"cgi", "jsp", "aspx", "EMPTY_STRING", "php", "py", "asp"}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAdvancedSettingsPrefetch/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_prefetch.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
