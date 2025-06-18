package interfaces

import "app/models"

type ICertificateService interface	{
	Obtain(cert *models.Certificate) error
}

type IServiceService interface	{
	FetchRancherServices(ws *models.Workspace) (interface{}, error)
	FetchPortainerServices(ws *models.Workspace) (interface{}, error)
}

type IWorkspaceService interface {
	GenerateTraefikConfig(ws *models.Workspace) (error)
}