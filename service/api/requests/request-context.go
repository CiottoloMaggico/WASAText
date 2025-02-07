// Package requests /*
package requests

import (
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// RequestContext is the context of the request, for request-dependent parameters
type RequestContext struct {
	// ReqUUID is the request unique ID
	ReqUUID    uuid.UUID
	IssuerUUID *string

	// Logger is a custom field logger for the request
	Logger logrus.FieldLogger
}
