# Event Storming — MVP-0 Logística Inversa Farmacéutica

**Bounded context:** `logistica-inversa-farmaceutica`
**Versión:** 0.1 (sesión big-picture — Architecture Discovery)
**Fecha:** 5 de octubre de 2026
**Autor:** Barto912 (owner de producto/dominio)
**Base:** `docs/domain/logistica-inversa.md` v0.2 · HC-CANNING-PROTO-006-LIF v1.1 · Issue #2
**Estado de madurez:** architecture discovery — sin código de producción hasta el Architecture Gate E1.

---

## 0. Convenciones de lectura

| Elemento | Notación | Significado |
|---|---|---|
| Evento de dominio | `PascalCase en pasado` | Algo que ocurrió y le importa al negocio |
| Comando | verbo imperativo | Acción que dispara un evento |
| Actor | **negrita** | Persona/organización con autoridad sobre el comando |
| Política | `P#` | Regla reactiva: "cuando X, entonces Y" |
| Sistema externo | `[corchetes]` | Fuera del boundary del contexto |
| Hotspot | `H#` | Incertidumbre, conflicto o pregunta abierta |

---

## 1. Big picture (línea de tiempo)

```
MedicamentoVencidoDetectado
  → ConsentimientoObtenido
    → ResiduoClasificado
      → RetiroConsolidadoEnBase
        → ManifiestoGenerado
          → TransferenciaATransportista
            → RecepcionEnPlanta
              → DisposicionFinalCertificada
                → RetiroCerrado

Ramas transversales:
  RecallANMATActivado ────┐
  CadenaDeFrioRota ───────┤→ FailCloseActivado → BloqueadoFailClose → BloqueoResuelto
  Invalidación (I6/I2/I4) ┘
  ConsentimientoRechazado → RetiroCerrado (motivo registrado)
```

---

## 2. Segmentos detallados

### Segmento 1 — Detección y consentimiento

| Evento | Comando | Actor | Política / nota |
|---|---|---|---|
| `MedicamentoVencidoDetectado` | DetectarMedicamento · LeerMedicationIdentity | **EnfermeroRetirador** (lectura: Sistema) | I6 como hipótesis (H4): MedicationIdentity completa o Fail-Close |
| — | SolicitarConsentimiento | Sistema (en nombre de Enfermero) | P1 |
| `ConsentimientoObtenido` | RegistrarConsentimiento | **Paciente/Socio** | Receipt operacional de consentimiento (H3) |
| `ConsentimientoRechazado` | RegistrarRechazo | **Paciente/Socio** | P2 |

### Segmento 2 — Clasificación y consolidación

| Evento | Comando | Actor | Política / nota |
|---|---|---|---|
| `ResiduoClasificado` | ClasificarResiduo | **FarmacéuticoClasificador** | PELIGROSO / ASIMILABLE (MVP-0) `[VALIDAR-FARMACIA]` |
| `RetiroConsolidadoEnBase` | ConsolidarEnBase | **EnfermeroRetirador** | Compartimento inverso sellado; telemetría activa |
| `CadenaDeFrioRota` (rama) | — (telemetría) | [IoT cadena de frío] | P3 |

### Segmento 3 — Manifiesto y primera transferencia

| Evento | Comando | Actor | Política / nota |
|---|---|---|---|
| `ManifiestoGenerado` | EmitirManifiesto | **FarmacéuticoClasificador** | I3 · hotspot H2 |
| `TransferenciaATransportista` | TransferirATransportista | **Farmacéutico → TransportistaHabilitado** | I2 firma dual · I4 habilitación vigente al momento (H6) · P4 |

### Segmento 4 — Planta y disposición final

| Evento | Comando | Actor | Política / nota |
|---|---|---|---|
| `RecepcionEnPlanta` | RecibirEnPlanta | **OperadorDisposicion** | I2 · I5 sin saltos |
| `DisposicionFinalCertificada` | CertificarDisposicionFinal | **OperadorDisposicion** | P5 · pregunta R4 |
| `RetiroCerrado` | CerrarRetiro | **FarmacéuticoClasificador** | I5 cadena completa |

### Segmento 5 — Fail-Close (transversal)

| Evento | Comando | Actor | Política / nota |
|---|---|---|---|
| `FailCloseActivado` | BloquearFailClose | SistemaGentleCare (automático) | I9 · triggers: I6 incompleta, firma faltante, habilitación vencida, identidad ilegible, CadenaDeFrioRota |
| `BloqueoResuelto` | ResolverBloqueo | **FarmacéuticoClasificador o Director Médico** | Nunca el Sistema |

---

## 3. Políticas (reglas reactivas)

| # | Cuando… | Entonces… |
|---|---|---|
| P1 | Medicamento retirado estaba activo | Notificar al MédicoTratante antes de continuar |
| P2 | ConsentimientoRechazado | Cierre con registro (RetiroCerrado con motivo) |
| P3 | CadenaDeFrioRota | Cuarentena + bloqueo de transferencia hasta validación farmacéutica |
| P4 | ManifiestoGenerado sin transferencia en 72 h | Escalada a Director Médico + transportista alternativo del registro vigente |
| P5 | RecepcionEnPlanta sin certificado en 30 días | Alerta regulatoria interna; evaluar reporte a autoridad `[VALIDAR-LEGAL R4]` |
| P6 | Cualquier invalidación de invariante | FailCloseActivado automático |
| P7 | RecallANMATActivado | Régimen específico de notificación/retiro `[VALIDAR-LEGAL R1]` |
| P8 | BloqueadoFailClose | Solo Farmacéutico o Director Médico pueden resolver |

---

## 4. Sistemas externos

| Sistema | Interacción | Frontera |
|---|---|---|
| [ANMAT — Trazabilidad Disp. 3683/2011] | Consulta GTIN+serie, recalls | Solo lectura |
| [OPDS — Registro transportistas/plantas PBA] | Validez de habilitaciones | Solo lectura + cache con refresh (H6) |
| [Google Cloud Healthcare API — FHIR] | Publicación de eventos clínicos relevantes | Escritura con consentimiento |
| [IoT cadena de frío] | Telemetría del compartimento inverso | Solo lectura |
| [Canal MédicoTratante] | Notificaciones P1 | Interno |

Ningún sistema externo posee `RetiroFarmaceutico`: el aggregate vive dentro del boundary del contexto.

---

## 5. Hotspots

| # | Hotspot | Tipo | Semilla |
|---|---|---|---|
| H1 | ¿`CadenaDeCustodia` es aggregate propio o colección de entities dentro de `RetiroFarmaceutico`? | Boundary | ADR-001 |
| H2 | ¿`Manifiesto` vive en este contexto o en un contexto regulatorio separado? | Boundary | ADR-001 |
| H3 | Receipt operacional: esquema de firma, timestamp confiable y anclaje para autenticidad y no-repudio (SHA-256 = solo integridad) | Técnico-legal | ADR-002 |
| H4 | I6 en terreno: ¿la MedicationIdentity completa es siempre legible? ¿Fallback definido si el soporte no trae serie? | Validación de campo | Farmacéutico + piloto |
| H5 | Recall ANMAT sin consentimiento: régimen aplicable | Legal | R1 |
| H6 | Validación temporal de habilitaciones (I4): fuente de verdad OPDS, refresh y comportamiento offline en borde | Técnico | ADR-003 |
| H7 | Tensión Ley 25.326 (derecho de cancelación/ARCO) vs. store append-only de receipts | Legal-técnico | ADR-002 |
| H8 | Timeouts 72 h / 30 días: ¿contractuales u operativos? ¿Dónde se codifica el SLA? | Operativo | Context Map |

---

## 6. Mapa hotspots → ADRs y validaciones

- **ADR-001 Aggregate boundaries** ← H1, H2
- **ADR-002 Operational evidence** (firma, timestamp, anclaje, ARCO vs inmutabilidad) ← H3, H7
- **ADR-003 Registros externos y operación offline** ← H6
- **Validación legal** ← H5 (R1), P5 (R4)
- **Validación farmacéutica / piloto de campo** ← H4

---

## 7. Notas de alineación para discovery v0.3

- `EventoCerrado` → `RetiroCerrado` y `CerrarEvento` → `CerrarRetiro` (consistencia con el rename del aggregate).
- Eventos nuevos detectados en la tormenta: `RetiroConsolidadoEnBase`, `BloqueoResuelto` (I8 exige un evento por transición; el discovery v0.2 no los nombraba).
- I6 **no se toca**: sigue como hipótesis hasta validación.

---

## 8. Trazabilidad

| Sección | Origen |
|---|---|
| §1-2 Timeline y segmentos | Discovery v0.2 §4-§6 |
| §3 Políticas | Discovery v0.2 §7-§8 · PROTO-006-LIF v1.1 §8 |
| §4 Sistemas externos | Discovery v0.2 §9 |
| §5 Hotspots | Discovery v0.2 §10-§12 · preguntas R1-R6 |

---

## 9. Próximos pasos del recorrido

1. Taller **Aggregate Boundaries** (resolver H1, H2)
2. **Context Map** (contextos vecinos: regulatorio, clínico-FHIR, logística IoT)
3. **ADRs 001-003**
4. **Architecture Gate E1**
