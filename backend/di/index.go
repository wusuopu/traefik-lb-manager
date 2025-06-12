package di

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type container struct {
	DB *gorm.DB
	Logger *zap.Logger
}

type service struct {
}

var Container = new(container)
var Service = new(service)

