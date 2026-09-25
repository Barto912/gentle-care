```
package cloud

import (
	"encoding/json"
	"fmt"
	"time"

	"gentle-care/pkg/triage"
)

// Resource FHIR R4 Observation para Google Cloud Healthcare API
type FHIRObservation struct {
	ResourceType string `json:"resourceType"` // "Observation"
	ID           string `json:"id"`
	Status       string `json:"status"` // "final"
	Code         struct {
		Coding []struct {
			System string `json:"system"` // "http://snomed.info/sct"
			Code   string `json:"code"`
			Display string `json:"display"`
		} `json:"coding"`
	} `json:"code"`
	Subject struct {
		Reference string `json:"reference"`
	} `json:"subject"`
	EffectiveDateTime string `json:"effectiveDateTime"`
	ValueString       string `json:"valueString"`
}

// BuildFHIRObservation transforma la evaluación de Gentle Care a estándar FHIR R4
func BuildFHIRObservation(assessment triage.Assessment, receiptSHA256 string) ([]byte, error) {
	obs := FHIRObservation{
		ResourceType: "Observation",
		ID:           assessment.ID,
		Status:       "final",
		EffectiveDateTime: assessment.Timestamp.Format(time.RFC3339),
		ValueString:       fmt.Sprintf("Triage Category: %s | Receipt SHA-256: %s", assessment.Category, receiptSHA256),
	}

	obs.Subject.Reference = "Patient/" + assessment.PatientIDHash
	obs.Code.Coding = []struct {
		System  string `json:"system"`
		Code    string `json:"code"`
		Display string `json:"display"`
	}{
		{
			System:  "http://snomed.info/sct",
			Code:    assessment.SnomedCode,
			Display: "Evaluación Domiciliaria Triage",
		},
	}

	return json.MarshalIndent(obs, "", "  ")
}

```
