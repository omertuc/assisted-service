package hostcommands

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/models"
	"github.com/sirupsen/logrus"
)

type domainResolutionCmd struct {
	baseCmd
	domainResolutionImage string
	db                    *gorm.DB
}

func NewDomainResolutionCmd(log logrus.FieldLogger, domainResolutionImage string, db *gorm.DB) *domainResolutionCmd {
	return &domainResolutionCmd{
		baseCmd:               baseCmd{log: log},
		domainResolutionImage: domainResolutionImage,
		db:                    db,
	}
}

func (d *domainResolutionCmd) prepareParam(host *models.Host, cluster *common.Cluster) (string, error) {
	request := models.DomainResolutionRequest{
		Domains: []*models.DomainResolutionRequestDomain{
			{DomainName: cluster.APIVipDNSName},
		},
	}
	b, err := json.Marshal(&request)
	if err != nil {
		d.log.WithError(err).Warn("Json marshal")
		return "", err
	}
	return string(b), nil
}

func (d *domainResolutionCmd) GetSteps(ctx context.Context, host *models.Host) ([]*models.Step, error) {
	var cluster common.Cluster
	if err := d.db.Take(&cluster, "id = ?", host.ClusterID.String()).Error; err != nil {
		return nil, err
	}

	param, err := d.prepareParam(host, &cluster)
	if err != nil {
		return nil, err
	}
	step := &models.Step{
		StepType: models.StepTypeDomainResolution,
		Command:  "podman",
		Args: []string{
			"run", "--privileged", "--net=host", "--rm", "--quiet",
			"-v", "/var/log:/var/log",
			"-v", "/run/systemd/journal/socket:/run/systemd/journal/socket",
			d.domainResolutionImage,
			"domain_resolution",
			param,
		},
	}

	return []*models.Step{step}, nil
}
