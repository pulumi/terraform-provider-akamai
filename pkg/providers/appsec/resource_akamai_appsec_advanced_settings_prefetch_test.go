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

		cd := appsec.UpdateAdvancedSettingsPrefetchResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"))
		json.Unmarshal([]byte(expectJSD), &cd)

		cr := appsec.GetAdvancedSettingsPrefetchResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetAdvancedSettingsPrefetch",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAdvancedSettingsPrefetchRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateAdvancedSettingsPrefetch",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAdvancedSettingsPrefetchRequest{ConfigID: 43253, Version: 7, AllExtensions: false, EnableAppLayer: true, EnableRateControls: false, Extensions: []string{"cgi", "asp", "php", "jsp", "EMPTY_STRING", "py", "aspx"}},
		).Return(&cu, nil)

		client.On("UpdateAdvancedSettingsPrefetch",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAdvancedSettingsPrefetchRequest{ConfigID: 43253, Version: 7, AllExtensions: false, EnableAppLayer: false, EnableRateControls: false, Extensions: []string(nil)},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
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
