package host

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-openapi/strfmt"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/models"
)

var _ = FDescribe("Majority groups", func() {
	Context("Test majority group validation", func() {
		var (
			cluster        common.Cluster
			majorityGroups map[string][]strfmt.UUID
			v              validator
		)

		BeforeEach(func() {
			cluster = common.Cluster{
				Cluster: models.Cluster{
					MachineNetworks: []*models.MachineNetwork{
						{
							Cidr:      "10.131.30.64/27",
							ClusterID: "12185c0e-fa62-48ab-a3c7-42a1725bdecf",
						},
						{
							Cidr:      "2001:1b74::/32",
							ClusterID: "12185c0e-fa62-48ab-a3c7-42a1725bdecf",
						},
					},
				},
			}

			majorityGroups = map[string][]strfmt.UUID{
				"10.131.30.64/27": {
					"1719d65d-ecc8-8e07-b868-284711f738ba",
					"a92a8904-7aac-0542-1289-c7b4260036f9",
					"f95dd0f7-e328-dd4f-5fdc-bd329ff77153",
				},
				"2001:1b74:480:6140::/64": {
					"1719d65d-ecc8-8e07-b868-284711f738ba",
					"a92a8904-7aac-0542-1289-c7b4260036f9",
					"f95dd0f7-e328-dd4f-5fdc-bd329ff77153",
				},
				"IPv4": {
					"1719d65d-ecc8-8e07-b868-284711f738ba",
					"a92a8904-7aac-0542-1289-c7b4260036f9",
					"f95dd0f7-e328-dd4f-5fdc-bd329ff77153",
				},
				"IPv6": {
					"1719d65d-ecc8-8e07-b868-284711f738ba",
					"a92a8904-7aac-0542-1289-c7b4260036f9",
					"f95dd0f7-e328-dd4f-5fdc-bd329ff77153",
				},
			}

			v = validator{}
		})

		It("Test belongs to majority group", func() {
			for _, id := range []strfmt.UUID{
				"1719d65d-ecc8-8e07-b868-284711f738ba",
				"a92a8904-7aac-0542-1289-c7b4260036f9",
				"f95dd0f7-e328-dd4f-5fdc-bd329ff77153",
			} {
				Expect(v.belongsToL2MajorityGroup(&validationContext{
					host: &models.Host{
						ID: &id,
					},
					cluster: &cluster,
				}, majorityGroups)).To(Equal(ValidationSuccess))
			}
		})
	})

})
