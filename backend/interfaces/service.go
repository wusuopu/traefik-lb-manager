package interfaces

import "app/models"

type ICertificateService interface	{
	Obtain(cert *models.Certificate) error
}