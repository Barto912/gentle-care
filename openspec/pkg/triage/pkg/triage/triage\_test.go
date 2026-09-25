```
package triage

import (
	"testing"
)

func TestEvaluateTriage(t *testing.T) {
	symptoms := []string{"Fiebre", "Disnea leve"}
	snomedCode := "267036007" // Código SNOMED CT para Disnea

	assessment, receipt := Evaluate(symptoms, false, snomedCode)

	if assessment.Category != CodeCeleste &amp;&amp; assessment.Category != CodeAmarillo {
		t.Errorf("Categoría inesperada: %s", assessment.Category)
	}

	if len(receipt) != 64 {
		t.Errorf("Firma SHA-256 inválida, longitud: %d", len(receipt))
	}
}

```
