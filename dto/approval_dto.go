package dto

// ApprovalRequest adalah body request untuk aksi approve/reject/send
type ApprovalRequest struct {
	Action  string `json:"action"  validate:"required"` // "approve" | "reject" | "send"
	Catatan string `json:"catatan"`
}

// SendWAResponse adalah data yang dikembalikan saat HRGA trigger kirim WA
type SendWAResponse struct {
	WAUrl       string `json:"wa_url"`        // URL click-to-chat WhatsApp
	NoSppd      string `json:"no_sppd"`
	NamaDireksi string `json:"nama_direksi"`
}
