package valueobject

type AdditionalRepairStatus string

const (
	StatusARAberta              AdditionalRepairStatus = "ABERTA"
	StatusARAguardandoAprovacao AdditionalRepairStatus = "AGUARDANDO_APROVACAO"
	StatusAAprovada             AdditionalRepairStatus = "APROVADA"
	StatusARRejeitada           AdditionalRepairStatus = "REJEITADA"
)

func ParseAdditionalRepairStatus(status string) AdditionalRepairStatus {
	switch status {
	case "ABERTA":
		return StatusARAberta
	case "AGUARDANDO_APROVACAO":
		return StatusARAguardandoAprovacao
	case "APROVADA":
		return StatusAAprovada
	case "REJEITADA":
		return StatusARRejeitada
	default:
		return AdditionalRepairStatus(status)
	}
}

func (s AdditionalRepairStatus) String() string {
	return string(s)
}
