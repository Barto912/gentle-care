package triage

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	symptoms := []string{"Fiebre"}
	assessment, receipt := Evaluate(symptoms, false, "267036007")

	if assessment.Category != CodeCeleste {
		t.Errorf("Error en categoria: %s", assessment.Category)
	}

	if len(receipt) != 64 {
		t.Errorf("Error en hash: %d", len(receipt))
	}
}
