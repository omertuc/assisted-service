package migrations

import (
	"encoding/json"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/onsi/gomega/types"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/models"
	gormigrate "gopkg.in/gormigrate.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func expectAllDiskEligibilityTo(diskEligibilityExpectation types.GomegaMatcher, db *gorm.DB) {
	hosts, err := db.Model(&models.Host{}).Rows()
	Expect(err).NotTo(HaveOccurred())

	var host models.Host
	for hosts.Next() {
		Expect(db.ScanRows(hosts, &host)).NotTo(HaveOccurred())

		var inventory models.Inventory
		Expect(json.Unmarshal([]byte(host.Inventory), &inventory)).NotTo(HaveOccurred())

		for _, disk := range inventory.Disks {
			Expect(disk.InstallationEligibility).To(diskEligibilityExpectation)
		}
	}
}

func createTestCluster(inventories ...string) *common.Cluster {
	someString := "arbitrary"

	clusterID := strfmt.UUID(uuid.New().String())

	cluster := &common.Cluster{Cluster: models.Cluster{
		ID:    &clusterID,
		Hosts: []*models.Host{},
	}}

	for _, inventory := range inventories {
		hostID := strfmt.UUID(uuid.New().String())
		cluster.Hosts = append(cluster.Hosts,
			&models.Host{
				Inventory: inventory, Kind: &someString, ID: &hostID,
				Href: &someString, Status: &someString, StatusInfo: &someString,
			})
	}

	return cluster
}

var _ = Describe("populateDiskEligibility", func() {
	var (
		db  *gorm.DB
		gm  *gormigrate.Gormigrate
		err error
	)

	AfterEach(func() {
		common.DeleteTestDB(db, "populate_disk_eligibility")
	})

	BeforeEach(func() {
		db = common.PrepareTestDB("populate_disk_eligibility")
		gm = gormigrate.New(db, gormigrate.DefaultOptions, all())
	})

	Context("General migration", func() {
		BeforeEach(func() {
			err = db.Create(createTestCluster(`{"disks": [{}, {}]}`)).Error
			Expect(err).NotTo(HaveOccurred())

			err = db.Create(createTestCluster(`{}`)).Error
			Expect(err).NotTo(HaveOccurred())

			err = db.Create(createTestCluster()).Error
			Expect(err).NotTo(HaveOccurred())

			err = gm.MigrateTo("20201215140033")
			Expect(err).ToNot(HaveOccurred())
		})

		It("Migrates down and up", func() {
			expectAllDiskEligibilityTo(Not(BeNil()), db)

			err = gm.RollbackMigration(populateDiskEligibility())
			Expect(err).ToNot(HaveOccurred())
			expectAllDiskEligibilityTo(BeNil(), db)

			err = gm.MigrateTo("20201215140033")
			Expect(err).ToNot(HaveOccurred())
			expectAllDiskEligibilityTo(Not(BeNil()), db)
		})
	})

	Context("Avoid corruption", func() {
		BeforeEach(func() {
			inventory := `
{
	"disks": [
		{
			"installation_eligibility": {
				"eligible":false,
				"not_eligible_reasons":["Disk is bad"]
			}
		}
	]
}`
			err = db.Create(createTestCluster(inventory)).Error
			Expect(err).NotTo(HaveOccurred())

			err = gm.MigrateTo("20201215140033")
			Expect(err).ToNot(HaveOccurred())
		})

		It("Doesn't corrupt existing entries", func() {
			expectAllDiskEligibilityTo(Equal(
				&models.DiskInstallationEligibility{
					Eligible:           false,
					NotEligibleReasons: []string{"Disk is bad"},
				},
			), db)
		})
	})

})
