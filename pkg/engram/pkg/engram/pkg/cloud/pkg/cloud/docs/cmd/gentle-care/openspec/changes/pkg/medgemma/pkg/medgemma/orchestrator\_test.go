```
ackage medgemma

import (
	"testing"
)

func TestInferTriageOfflineFallback(t *testing.T) {
	orchestrator := NewOrchestrator(false) // Modo Offline / Sin Conexión

	assessment, receipt, mode := orchestrator.InferTriage([]string{"Fiebre"}, "267036007")

	if mode != "LOCAL_EDGE_DETERMINISTIC_FAIL_CLOSE" {
		t.Errorf("Se esperaba conmutación a contingencia local determinista, obtenido: %s", mode)
	}

	if assessment.Category == "" || len(receipt) != 64 {
		t.Errorf("Evaluación o recibo SHA-256 inválido durante contingencia offline")
	}
}

```
