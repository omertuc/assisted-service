package featuresupport

import (
	"github.com/openshift/assisted-service/internal/usage"
	"github.com/openshift/assisted-service/models"
	"sort"
)

func usageNameToID(key string) string {
	return usage.UsageNameToID(key)
}

var SupportLevelsListWithUnsupported models.FeatureSupportLevels

// GetAllFeatureIds returns all feature IDs that appear somewhere in
// SupportLevelsListConcise
func GetAllFeatureIds() []string {
	featureIdSet := map[string]bool{}

	for _, versionFeatures := range SupportLevelsListConcise {
		for _, feature := range versionFeatures.Features {
			featureIdSet[feature.FeatureID] = true
		}
	}

	allFeatureIds := []string{}
	for featureId := range featureIdSet {
		allFeatureIds = append(allFeatureIds, featureId)
	}

	sort.Strings(allFeatureIds)

	return allFeatureIds
}

// PopulateUnsupportedFeatures fills SupportLevelsListVerbose with entries that
// say "Unsupported" for every feature that was not explicitly mentioned for a
// particular version in SupportLevelsListConcise. This allows us to not be
// forced to list every newly added feature as unsupported in a hard-coded
// manner for all historical versions, while still allowing us to return a
// fully populated list when requested in the API.
// This function is called on init.
func PopulateUnsupportedFeatures() {
	allFeatureIds := GetAllFeatureIds()

	for _, versionFeatures := range SupportLevelsListConcise {
		for _, feature := range allFeatureIds {
			versionFeatures.Features = append(versionFeatures.Features, &models.FeatureSupportLevelFeaturesItems0{
				FeatureID:    feature,
				SupportLevel: GetFeatureSupportLevel(versionFeatures.OpenshiftVersion, feature),
			})
		}
	}
}

func init() {
	PopulateUnsupportedFeatures()
}

var SupportLevelsListVerbose models.FeatureSupportLevels

// A mapping between versions and their support level of each feature.
// Do not add explicit "unsupported" entries for features, as features
// are assumed to be unsupported when absent from this list.
var SupportLevelsListConcise = models.FeatureSupportLevels{
	&models.FeatureSupportLevel{
		OpenshiftVersion: "4.6",
		Features: []*models.FeatureSupportLevelFeaturesItems0{
			// Dev-Preview features
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDVIPAUTOALLOC,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelDevPreview,
			},
		},
	},
	&models.FeatureSupportLevel{
		OpenshiftVersion: "4.8",
		Features: []*models.FeatureSupportLevelFeaturesItems0{
			// Dev-Preview features
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDVIPAUTOALLOC,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelDevPreview,
			},
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDSNO,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelDevPreview,
			},
		},
	},
	&models.FeatureSupportLevel{
		OpenshiftVersion: "4.9",
		Features: []*models.FeatureSupportLevelFeaturesItems0{
			// Supported
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDSNO,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelSupported,
			},
			// Dev-Preview features
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDVIPAUTOALLOC,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelDevPreview,
			},
		},
	},
	&models.FeatureSupportLevel{
		OpenshiftVersion: "4.10",
		Features: []*models.FeatureSupportLevelFeaturesItems0{
			// Supported
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDSNO,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelSupported,
			},
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDARM64ARCHITECTURE,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelSupported,
			},
			// Dev-Preview features
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDVIPAUTOALLOC,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelDevPreview,
			},
		},
	},
	&models.FeatureSupportLevel{
		OpenshiftVersion: "4.11",
		Features: []*models.FeatureSupportLevelFeaturesItems0{
			// Supported
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDSNO,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelSupported,
			},
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDARM64ARCHITECTURE,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelSupported,
			},
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDARM64ARCHITECTUREWITHCLUSTERMANAGEDNETWORKING,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelSupported,
			},
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDSINGLENODEEXPANSION,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelSupported,
			},
			// Dev-Preview features
			{
				FeatureID:    models.FeatureSupportLevelFeaturesItems0FeatureIDVIPAUTOALLOC,
				SupportLevel: models.FeatureSupportLevelFeaturesItems0SupportLevelDevPreview,
			},
		},
	},
}

// GetFeatureSupportLevel consults the list of supported features for a
// particular OCP release version (SupportLevelsListConcise) to decide on the
// support level of a particular feature. If the support level for a particular
// feature on the requested release version is not explicitly specified, it is
// assumed to be unsupported.
func GetFeatureSupportLevel(openshiftVersion string, featureId string) string {
	for _, supportLevel := range SupportLevelsListConcise {
		if supportLevel.OpenshiftVersion == openshiftVersion {
			for _, feature := range supportLevel.Features {
				if usageNameToID(featureId) == feature.FeatureID {
					return feature.SupportLevel
				}
			}
			break
		}
	}
	return models.FeatureSupportLevelFeaturesItems0SupportLevelUnsupported
}

func IsFeatureSupported(openshiftVersion string, featureId string) bool {
	return GetFeatureSupportLevel(openshiftVersion, featureId) == models.FeatureSupportLevelFeaturesItems0SupportLevelSupported
}
