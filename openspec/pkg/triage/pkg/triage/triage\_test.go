package triage

import (
	"testing"
)

func TestEvaluateTriage(t *testing.T) {
	symptoms := []string{"Fiebre", "Disnea leve"}
	snomedCode := "267036007"

	assessment, receipt := Evaluate(symptoms, false, snomedCode)

	if assessment.Category != CodeCeleste && assessment.Category != CodeAmarillo {
		t.Errorf("Categoria inesperada: %s", assessment.Category)
	}

	if len(receipt) != 64 {
		t.Errorf("Firma SHA-256 invalida, longitud: %d", len(receipt))
	}
}
