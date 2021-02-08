module github.com/openshift/assisted-service/client

go 1.15

replace github.com/openshift/assisted-service/models => ../models

require (
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/runtime v0.19.20
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/openshift/assisted-service/models v0.0.0-00010101000000-000000000000
)

replace github.com/openshift/assisted-service/models => ../models

