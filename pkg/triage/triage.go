package triage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type TriageCode string

const (
	CodeCeleste  TriageCode = "CELESTE"
	CodeAmarillo TriageCode = "AMARILLO"
	CodeRojo     TriageCode = "ROJO"
)

type Assessment struct {
	ID            string     `json:"id"`
	PatientIDHash string     `json:"patient_id_hash"`
	SnomedCode    string     `json:"snomed_code"`
	Category      TriageCode `json:"category"`
	Symptoms      []string   `json:"symptoms"`
	Timestamp     time.Time  `json:"timestamp"`
}

func Evaluate(symptoms []string, isCritical bool, snomed string) (Assessment, string) {
	var category TriageCode

	if isCritical {
		category = CodeRojo
	} else if len(symptoms) > 2 {
		category = CodeAmarillo
	} else {
		category = CodeCeleste
	}

	eval := Assessment{
		ID:            fmt.Sprintf("EVAL-%d", time.Now().UnixNano()),
		PatientIDHash: "ANON-PHI-ENCODED",
		SnomedCode:    snomed,
		Category:      category,
		Symptoms:      symptoms,
		Timestamp:     time.Now().UTC(),
	}

	hashInput := fmt.Sprintf("%s|%s|%s|%s", eval.ID, eval.SnomedCode, eval.Category, eval.Timestamp.String())
	hash := sha256.Sum256([]byte(hashInput))
	receiptSHA256 := hex.EncodeToString(hash[:])

	return eval, receiptSHA256
}
