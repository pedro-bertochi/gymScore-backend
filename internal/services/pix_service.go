package services

import (
	"fmt"
	"gynScore-backend/internal/config"
	"gynScore-backend/internal/models"
	"gynScore-backend/internal/repositories"
	"gynScore-backend/pkg/utils"
)

type PIXService interface {
	GerarPagamento(req models.PIXRequest) (*models.PIXResponse, error)
}

type pixService struct {
	cfg         *config.Config
	usuarioRepo repositories.UsuarioRepository
}

func NovoPIXService(cfg *config.Config, usuarioRepo repositories.UsuarioRepository) PIXService {
	return &pixService{cfg, usuarioRepo}
}

func (s *pixService) GerarPagamento(req models.PIXRequest) (*models.PIXResponse, error) {
	// Validações
	if req.Valor <= 0 {
		return nil, fmt.Errorf("valor do depósito deve ser maior que zero")
	}

	if !utils.ValidarCPF(req.CPF) {
		return nil, fmt.Errorf("CPF informado é inválido")
	}

	usuario, err := s.usuarioRepo.BuscarPorID(req.IDUsuario)
	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado: %v", err)
	}

	// 🔥 TXID seguro (sem caracteres problemáticos)
	descricao := fmt.Sprintf("DEP%d", usuario.ID)

	payload := utils.GeneratePIXPayload(
		s.cfg.PIXChave,
		s.cfg.PIXNome,
		s.cfg.PIXCidade,
		descricao,
		req.Valor,
	)

	qrcodeBase64, err := utils.GenerateQRCodeBase64(payload)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar QR Code: %v", err)
	}

	// ⚠️ SIMULAÇÃO (não é pagamento real)
	usuario.Saldo += req.Valor
	if err := s.usuarioRepo.Atualizar(usuario); err != nil {
		return nil, fmt.Errorf("falha ao atualizar saldo: %v", err)
	}

	return &models.PIXResponse{
		QRCodeBase64: qrcodeBase64,
		Payload:      payload,
		NovoSaldo:    usuario.Saldo,
	}, nil
}
