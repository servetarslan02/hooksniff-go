package models

type ApplicationAlertEvent string
const (
	ApplicationAlertEventEndpointDisabled ApplicationAlertEvent = "endpoint.disabled"
	ApplicationAlertEventEndpointEnabled  ApplicationAlertEvent = "endpoint.enabled"
)
