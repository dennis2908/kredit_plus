package structs

import(
	"kredit_plus/models"
)

type GetKonsumen struct {
	models.Konsumen
	Date_birth	 string
}

type InsertKonsumen struct {
	Konsumen models.Konsumen
	Date_birth	 string
}
type InsertKonsumenM struct {
	models.Konsumen
	Date_birth	 string
}