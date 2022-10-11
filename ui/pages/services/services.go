package services

import (
	"github.com/danielcmessias/lfq/ui/components/page"
)

type ServicesPageModel struct {
    page.Model
}

func NewServicesPage() *ServicesPageModel {
	return &ServicesPageModel{
		Model: page.New(servicesPageSpec),
	}
}
