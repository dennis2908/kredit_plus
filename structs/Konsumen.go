package structs

import (
	"kredit_plus/models"
)

type GetKonsumen struct {
	models.Konsumen
	Date_birth string
}
