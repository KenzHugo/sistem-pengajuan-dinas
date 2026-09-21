package dto

// TarifRequest untuk menambah / memperbarui data tarif perjalanandinas_tarif dengan backward compatibility
type TarifRequest struct {
	NoUrut          string  `json:"no_urut"`
	Area            string  `json:"area"`
	AreaTujuan      string  `json:"area_tujuan"`
	Jabatan         string  `json:"jabatan"`
	KriteriaJabatan string  `json:"kriteria_jabatan"`
	MataUang        string  `json:"mata_uang"`
	UangMakan       float64 `json:"uang_makan"`
	TarifUangMakan  float64 `json:"tarif_uang_makan"`
	UangSaku        float64 `json:"uang_saku"`
	TarifUangSaku   float64 `json:"tarif_uang_saku"`
	Hotel           float64 `json:"hotel"`
	HakHotel        float64 `json:"hak_hotel"`
	Transport       float64 `json:"transport"`
	HakTransport    float64 `json:"hak_transport"`
	Tiket           float64 `json:"tiket"`
}

// GetArea mengembalikan area tujuan dari field mana pun yang terisi
func (t *TarifRequest) GetArea() string {
	if t.Area != "" {
		return t.Area
	}
	return t.AreaTujuan
}

// GetJabatan mengembalikan nama jabatan / kriteria dari field mana pun yang terisi
func (t *TarifRequest) GetJabatan() string {
	if t.Jabatan != "" {
		return t.Jabatan
	}
	return t.KriteriaJabatan
}

// GetUangMakan mengembalikan nominal uang makan
func (t *TarifRequest) GetUangMakan() float64 {
	if t.UangMakan > 0 {
		return t.UangMakan
	}
	return t.TarifUangMakan
}

// GetUangSaku mengembalikan nominal uang saku
func (t *TarifRequest) GetUangSaku() float64 {
	if t.UangSaku > 0 {
		return t.UangSaku
	}
	return t.TarifUangSaku
}

// GetHotel mengembalikan nominal hotel
func (t *TarifRequest) GetHotel() float64 {
	if t.Hotel > 0 {
		return t.Hotel
	}
	return t.HakHotel
}

// GetTransport mengembalikan nominal transport
func (t *TarifRequest) GetTransport() float64 {
	if t.Transport > 0 {
		return t.Transport
	}
	return t.HakTransport
}

// CreateTarifRequest alias untuk kompatibilitas
type CreateTarifRequest = TarifRequest

// UpdateTarifRequest alias untuk kompatibilitas
type UpdateTarifRequest = TarifRequest
