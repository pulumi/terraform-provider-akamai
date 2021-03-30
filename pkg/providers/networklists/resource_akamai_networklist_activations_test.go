package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiActivations_res_basic(t *testing.T) {
	t.Run("match by Activations ID", func(t *testing.T) {
		client := &mocknetworklists{}

		cu := networklists.RemoveActivationsResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResActivations/ActivationsDelete.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		ga := networklists.GetActivationsResponse{}
		expectJSR := compactJSON(loadFixtureBytes("testdata/TestResActivations/Activations.json"))
		json.Unmarshal([]byte(expectJSR), &ga)

		cr := networklists.CreateActivationsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResActivations/Activations.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetActivations",
			mock.Anything, // ctx is irrelevant for this test
			networklists.GetActivationsRequest{ActivationID: 547694},
		).Return(&ga, nil)

		client.On("CreateActivations",
			mock.Anything, // ctx is irrelevant for this test
			networklists.CreateActivationsRequest{Action: "ACTIVATE", Network: "STAGING", Comments: "", NotificationRecipients: []string{"martin@email.io"}},
		).Return(&cr, nil)

		client.On("RemoveActivations",
			mock.Anything, // ctx is irrelevant for this test
			networklists.RemoveActivationsRequest{ActivationID: 547694, Action: "DEACTIVATE", Network: "STAGING", Comments: "", NotificationRecipients: []string{"martin@email.io"}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "id", "547694"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
