package service

import (
	"os"
	"strings"
	"yorgut/internal/transport"
)

// InfoService tells about this Yorgut instance, from its configuration, so
// whoever runs one sets their own details instead of editing the frontend.
type InfoService struct{}

func NewInfoService() *InfoService {
	return &InfoService{}
}

// FindInstanceInfo returns the administrator's e-mail (ADMIN_EMAIL), nil when
// it isn't set.
func (i *InfoService) FindInstanceInfo() *transport.InstanceInfoResponse {
	response := transport.InstanceInfoResponse{}
	if email := strings.TrimSpace(os.Getenv("ADMIN_EMAIL")); email != "" {
		response.AdminEmail = &email
	}
	return &response
}
