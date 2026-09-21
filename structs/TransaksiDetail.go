package structs

import (
	models "kredit_plus/models"
)

type GetTransaksiDetails struct {
	Transaksis        []GetAllTransaksi
	Transaksi_details []models.Transaksi_details
}
