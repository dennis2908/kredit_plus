package structs

import (
	"kredit_plus/models"
)

type GetAllTransaksi struct {
	Data_transaksi models.TransaksiWIdTrans
	Data_konsumen  models.KonsumensWIdKons
}
