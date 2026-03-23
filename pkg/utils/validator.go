package utils

import (
	"regexp"
	"strings"
)

// emailRegex é a expressão regular para validação de e-mail
// Equivalente ao regex utilizado no servidor Java original
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,6}$`)

// ValidarEmail verifica se o e-mail fornecido possui formato válido
func ValidarEmail(email string) bool {
	email = strings.TrimSpace(email)
	return emailRegex.MatchString(email)
}

// ValidarSaldo verifica se o usuário possui saldo suficiente para participar do desafio
// Equivalente à lógica do endpoint /validar-saldo do servidor Java original
func ValidarSaldo(saldoUsuario, valorDesafio float64) bool {
	return saldoUsuario >= valorDesafio
}

// CalcularSaldosAposDesafio calcula os novos saldos após o encerramento de um desafio
// Equivalente à lógica do endpoint /finalizar-desafio do servidor Java original:
// - Vencedor recebe o valor apostado integralmente
// - Perdedor perde metade do valor apostado
func CalcularSaldosAposDesafio(saldoVencedor, saldoPerdedor, valorApostado float64) (novoSaldoVencedor, novoSaldoPerdedor float64) {
	novoSaldoVencedor = saldoVencedor + valorApostado
	novoSaldoPerdedor = saldoPerdedor - (valorApostado / 2)
	return
}
