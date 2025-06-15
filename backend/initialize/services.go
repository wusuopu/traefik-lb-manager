package initialize

import (
	"app/di"
	"app/interfaces"
	"app/services"
)

func InitServices() {
	di.Service.CertificateService = interfaces.ICertificateService(new(services.CertificateService))
}