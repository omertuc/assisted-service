package network

import (
	"encoding/base64"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/sirupsen/logrus"
)

const unmanagedResolvConf = `
[main]
rc-manager=unmanaged
`

const snoDnsmasqConfTemplate = `
address=/apps.{{.CLUSTER_NAME}}.{{.DNS_DOMAIN}}/{{.HOST_IP}}
address=/api-int.{{.CLUSTER_NAME}}.{{.DNS_DOMAIN}}/{{.HOST_IP}}
address=/api.{{.CLUSTER_NAME}}.{{.DNS_DOMAIN}}/{{.HOST_IP}}
`

const forceDnsmasqDnsNmDispatcherScriptTemplate = `
export IP="{{.HOST_IP}}"
export BASE_RESOLV_CONF=/run/NetworkManager/resolv.conf
if [ "$2" = "dhcp4-change" ] || [ "$2" = "dhcp6-change" ] || [ "$2" = "up" ] || [ "$2" = "connectivity-change" ]; then
    if ! grep -q "$IP" /etc/resolv.conf; then
      export TMP_FILE=$(mktemp /etc/forcedns_resolv.conf.XXXXXX)
      cp  $BASE_RESOLV_CONF $TMP_FILE
      chmod --reference=$BASE_RESOLV_CONF $TMP_FILE
      sed -i -e "s/{{.CLUSTER_NAME}}.{{.DNS_DOMAIN}}//" \
      -e "s/search /& {{.CLUSTER_NAME}}.{{.DNS_DOMAIN}} /" \
      -e "0,/nameserver/s/nameserver/& $IP\n&/" $TMP_FILE
      mv $TMP_FILE /etc/resolv.conf
    fi
fi
`

const snoDnsmasqIngnitionFilesJsonTemplate = `
[
  {
    "contents": {
      "source": "data:text/plain;charset=utf-8;base64,{{.DNSMASQ_CONFIG}}",
      "verification": {}
    },
    "filesystem": "root",
    "mode": 420,
    "path": "/etc/dnsmasq.d/single-node.conf"
  },
  {
    "contents": {
      "source": "data:text/plain;charset=utf-8;base64,{{.FORCE_DNS_SCRIPT}}",
      "verification": {}
    },
    "filesystem": "root",
    "mode": 365,
    "path": "/etc/NetworkManager/dispatcher.d/forcedns"
  },
  {
    "contents": {
      "source": "data:text/plain;charset=utf-8;base64,{{.UNMANAGED_RESOLV_CONF}}",
      "verification": {}
    },
    "filesystem": "root",
    "mode": 420,
    "path": "/etc/NetworkManager/conf.d/single-node.conf"
  }
]
`

const snoDnsmasqMachineConfigManifestTemplate = `
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: master
  name: 50-master-dnsmasq-configuration
spec:
  config:
    ignition:
      config: {}
      security:
        tls: {}
      timeouts: {}
      version: 2.2.0
    networkd: {}
    passwd: {}
    storage:
      files: {{.DNSMASQ_FILES_JSON}}
    systemd:
      units:
        - name: dnsmasq.service
          enabled: true
          contents: |
            [Unit]
            Description=Run dnsmasq to provide local dns for Single Node OpenShift
            Before=kubelet.service crio.service
            After=network.target

            [Service]
            ExecStart=/usr/sbin/dnsmasq -k

            [Install]
            WantedBy=multi-user.target
`

func GetDnsmasqIgnitionFilesJson(log logrus.FieldLogger, cluster *common.Cluster) (string, error) {
	hostIp, err := GetIpForSingleNodeInstallation(cluster, log)
	if err != nil {
		return "", err
	}

	dnsTemplateParams := map[string]interface{}{
		"CLUSTER_NAME": cluster.Cluster.Name,
		"DNS_DOMAIN":   cluster.Cluster.BaseDNSDomain,
		"HOST_IP":      hostIp,
	}

	log.Infof("Creating dnsmasq manifest with values: cluster name: %q, domain - %q, host ip - %q",
		cluster.Cluster.Name, cluster.Cluster.BaseDNSDomain, hostIp)

	dnsmasqConf, err := fillTemplate(dnsTemplateParams, snoDnsmasqConfTemplate, log)
	if err != nil {
		return "", err
	}

	dispatcherScript, err := fillTemplate(dnsTemplateParams, forceDnsmasqDnsNmDispatcherScriptTemplate, log)
	if err != nil {
		return "", err
	}

	ignitionFileContentsTemplateParams := map[string]interface{}{
		"DNSMASQ_CONFIG":        base64.StdEncoding.EncodeToString(dnsmasqConf),
		"FORCE_DNS_SCRIPT":      base64.StdEncoding.EncodeToString(dispatcherScript),
		"UNMANAGED_RESOLV_CONF": base64.StdEncoding.EncodeToString([]byte(unmanagedResolvConf)),
	}

	ignitionFilesJson, err := fillTemplate(ignitionFileContentsTemplateParams, snoDnsmasqIngnitionFilesJsonTemplate, log)
	if err != nil {
		return "", err
	}

	return string(ignitionFilesJson), nil
}

func createDnsmasqMachineConfigForSingleNode(log logrus.FieldLogger, cluster *common.Cluster) ([]byte, error) {
	ignitionFilesJson, err := GetDnsmasqIgnitionFilesJson(log, cluster)
	if err != nil {
		return nil, err
	}

	machineConfigTemplateParams := map[string]interface{}{
		"DNSMASQ_FILES_JSON": ignitionFilesJson,
	}

	machineConfig, err := fillTemplate(machineConfigTemplateParams, snoDnsmasqMachineConfigManifestTemplate, log)
	if err != nil {
		return nil, err
	}

	return machineConfig, nil
}
