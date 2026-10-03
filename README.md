# Gentle Care — Nurse Canning 24

Plataforma de enfermería de precisión domiciliaria para los countries y
barrios privados de Canning (Esteban Echeverría, Provincia de Buenos Aires).
Alianza Clínica Canning Health × socio tecnológico.

## Estado del proyecto

| Módulo | Estado de madurez |
|---|---|
| `pkg/triage` | implementation — tests deterministas en CI |
| `pkg/rdd` | design — VerificationReceipt SHA-256 (próximo) |
| `pkg/engram` | research — memoria local SQLite+FTS5 |
| Protocolos clínicos | design — HC-CANNING-PROTO-001 y 006-LIF en validación |

## Arquitectura (resumen)

- **Borde:** Google Pixel 11 Pro, inferencia local, sanitización PHI en el dispositivo.
- **Nube:** Google Cloud Healthcare API (FHIR + SNOMED CT), transporte BeyondCorp Zero Trust.
- **Gobernanza:** ODD + SDD + RDD, compuerta Fail-Close, recibos SHA-256.

## Desarrollo

- Go 1.23, sin dependencias externas en esta fase.
- CI en cada push a main: `go test ./...`, `go vet ./...`, `go build ./...`.
- Estados de madurez: research → design → prototype → implementation → validated.

## Licencia

Ver [LICENSE](LICENSE).
