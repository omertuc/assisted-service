package constants

const Kubeconfig = "kubeconfig"
const KubeconfigNoIngress = "kubeconfig-noingress"

// an arbitrary subdomain of apps.<cluster-name>.<base_domain> used by DNS
// validations to verify that *.apps wildcard is configured properly
const AppsSubDomainNameHostDNSValidation = "console-openshift-console"

// Standard cluster API subdomains
const APIClusterSubdomain = "api"
const InternalAPIClusterSubdomain = "api-int"

// non-existing domain name under clusterName.baseDomain for wildcard configuration check
const DNSWildcardFalseDomainName = "validateNoWildcardDNS"
