package di

import (
	"app/interfaces"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type container struct {
	DB *gorm.DB
	Logger *zap.Logger
}

type service struct {
	CertificateService interfaces.ICertificateService
}

var Container = new(container)
var Service = new(service)

