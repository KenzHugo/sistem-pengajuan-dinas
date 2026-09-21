package dto

// ApprovalRequest adalah body request untuk aksi approve/reject/send beserta koreksi biaya (opsional)
type ApprovalRequest struct {
	Action         string   `json:"action"  validate:"required"` // "approve" | "reject" | "send"
	Catatan        string   `json:"catatan"`
	BiayaMakan     *float64 `json:"biaya_makan,omitempty"`
	BiayaTransport *float64 `json:"biaya_transport,omitempty"`
	BiayaUangSaku  *float64 `json:"biaya_uang_saku,omitempty"`
	BiayaTiket     *float64 `json:"biaya_tiket,omitempty"`
	BiayaHotel     *float64 `json:"biaya_hotel,omitempty"`
	EstimasiBiaya  *float64 `json:"estimasi_biaya,omitempty"`
}

// SendWAResponse adalah data yang dikembalikan saat HRGA trigger kirim WA
type SendWAResponse struct {
	WAUrl       string `json:"wa_url"`        // URL click-to-chat WhatsApp
	NoSppd      string `json:"no_sppd"`
	NamaDireksi string `json:"nama_direksi"`
}
