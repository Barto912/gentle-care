package main

import (
	"fmt"
	"gentle-care/pkg/triage"
)

func main() {
	symptoms := []string{"Fiebre", "Disnea leve"}
	assessment, receipt := triage.Evaluate(symptoms, false, "267036007")
	fmt.Printf("Gentle Care Engine Started. Triage: %s | Receipt: %s\n", assessment.Category, receipt)
}
