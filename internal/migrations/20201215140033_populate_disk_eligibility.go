package migrations

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/openshift/assisted-service/models"
	"gopkg.in/gormigrate.v1"
)

// migrateHostInventory behaves similarly to (m *host.Manager) populateDisksEligibility in the sense that
// it sets all nil eligibility structs to "eligible".
func migrateHostInventory(inventoryString string, newValue *models.DiskInstallationEligibility) (string, error) {
	var inventory models.Inventory
	if err := json.Unmarshal([]byte(inventoryString), &inventory); err != nil {
		return "", err
	}

	for _, disk := range inventory.Disks {
		// Don't migrate if an eligibility struct is already there. Unless we do a rollback.
		if disk.InstallationEligibility != nil && newValue != nil {
			continue
		}

		disk.InstallationEligibility = newValue
	}

	result, err := json.Marshal(inventory)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func migrateHosts(tx *gorm.DB, newValue *models.DiskInstallationEligibility) error {
	hosts, err := tx.Model(&models.Host{}).Rows()

	if err != nil {
		return err
	}

	var host models.Host
	for hosts.Next() {
		if err = tx.ScanRows(hosts, &host); err != nil {
			return err
		}

		var err error
		host.Inventory, err = migrateHostInventory(host.Inventory, newValue)
		if err != nil {
			return err
		}

		if err := tx.Model(&host).Update("inventory", host.Inventory).Error; err != nil {
			return err
		}
	}

	return nil
}

func populateDiskEligibility() *gormigrate.Migration {
	migrate := func(tx *gorm.DB) error {
		return migrateHosts(tx,
			&models.DiskInstallationEligibility{
				Eligible:           true,
				NotEligibleReasons: make([]string, 0),
			})
	}

	rollback := func(tx *gorm.DB) error {
		return migrateHosts(tx, nil)
	}

	return &gormigrate.Migration{
		ID:       "20201215140033",
		Migrate:  migrate,
		Rollback: rollback,
	}
}
