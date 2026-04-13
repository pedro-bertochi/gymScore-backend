package models

// PIXRequest representa os dados necessários para gerar um pagamento PIX de depósito
type PIXRequest struct {
	IDUsuario uint    `json:"id_usuario" example:"1"`
	Valor     float64 `json:"valor" example:"50.00"`
	CPF       string  `json:"cpf" example:"123.456.789-00"`
}

// PIXResponse contém o QR Code e a linha digitável do PIX
type PIXResponse struct {
	QRCodeBase64 string  `json:"qrcode_base64" description:"QR Code em formato Base64 (PNG)"`
	Payload      string  `json:"payload" description:"Código PIX Copia e Cola"`
	NovoSaldo    float64 `json:"novo_saldo" description:"Saldo atualizado do usuário"`
}
