# Domain Discovery — Logística Inversa Farmacéutica

**Bounded context:** `logistica-inversa-farmaceutica`
**Versión:** 0.2 (borrador de discovery — Gate E0: PASS WITH CONDITIONS)
**Cambios v0.2:** (a) rename EventoDeRetiro → RetiroFarmaceutico según revisión del Gate E0 (commit e0745ac); (b) I6 recategorizada como hipótesis pendiente de validación farmacéutica.
**Fecha:** 4 de octubre de 2026
**Autor:** Barto912 (owner de producto/dominio)
**Revisores pendientes:** Farmacéutico Institucional · Asesor legal ambiental-sanitario · Revisión técnica (Engineering Gate)
**Referencias:** HC-CANNING-PROTO-006-LIF v1.1 · Issue #2 · PROTO-001 v1.0
**Estado de madurez:** discovery — de este documento **no se deriva código de producción** hasta aprobar el Gate.

---

## 0. Cómo leer este documento

- Secciones 1-9: consenso de discovery (lenguaje, ciclo de vida, reglas).
- Sección 10: semántica de receipts — **decisión técnica abierta**, no implementar antes del Gate.
- Secciones 11-12: preguntas abiertas e hipótesis explícitas.
- Los términos en `Código` son candidatos de lenguaje ubicuo, **no nombres de clases**.
- Etiquetas `[VALIDAR-FARMACIA]` y `[VALIDAR-LEGAL]` marcan puntos que requieren firma profesional.

---

## 1. Contexto y alcance MVP-0

Alianza Clínica Canning Health × socio tecnológico (Nurse Canning 24): retiro de medicamentos del domicilio de socios de 57 countries y barrios privados de Canning (Esteban Echeverría, PBA) y devolución al circuito regulado con cadena de custodia trazable.

**Dentro del MVP-0:** vencidos o por vencer (<60 días) · discontinuados por cambio terapéutico · cadena de frío comprometida · recall ANMAT.
**Fuera del MVP-0:** patogénicos, estupefacientes/psicotrópicos, citostáticos, radiactivos (PROTO-006-LIF v1.1 §2.5).

**Criterio de trazabilidad:** cada término, estado y regla de este documento debe rastrearse a una sección del PROTO-006-LIF v1.1 o a una pregunta regulatoria abierta (§11).

---

## 2. Lenguaje ubicuo

| Término | Definición | NO es / NO confundir con |
|---|---|---|
| `MedicationIdentity` | Identidad normativa de la unidad: **GTIN + serie + lote + vencimiento** (Disp. ANMAT 3683/2011). Value object. | "GTIN-13", "el código": el GTIN es un componente, no la identidad |
| `MecanismoDeLectura` | Infraestructura física que captura la MedicationIdentity: barcode 1D, DataMatrix, RFID. **Infraestructura, no dominio.** | "Scanner" como concepto de dominio |
| `RetiroFarmaceutico` | Aggregate root: instancia de retiro de una o más unidades desde un domicilio, desde la detección hasta la disposición final certificada. | "Ticket", "caso", "retiro" (ambiguo) |
| `CadenaDeCustodia` | Secuencia de `NodoDeTransferencia` que documenta cada cambio de posesión física. | "Trazabilidad": la trazabilidad es la propiedad; la cadena es la estructura |
| `NodoDeTransferencia` | Entity: cambio de posesión entre dos actores, con lugar, tiempo y firma dual. | "Entrega", "handover" |
| `Manifiesto` | Documento regulatorio de transporte (triplicado) que vincula el RetiroFarmaceutico con un transportista habilitado. | "Remito", "guía" |
| `OperadorHabilitado` | Aggregate de datos maestros: transportista o planta con habilitación OPDS/autoridad **vigente**. | "Vendor", "proveedor" |
| `ReceiptOperacional` | Evidencia firmada de un **hecho del dominio clínico-operativo** (consentimiento, transferencia, disposición). Ver §10. | RDD / DevelopmentReceipt |
| `DevelopmentReceipt` (RDD) | Evidencia del **proceso de desarrollo** de Gentle-AI (builds, gates, evaluaciones). Otro contexto, otro lifecycle. | ReceiptOperacional |
| `ClasificacionResiduo` | Decisión del farmacéutico: PELIGROSO / ASIMILABLE en MVP-0. `[VALIDAR-FARMACIA]` | "Diagnóstico" |
| `Fail-Close` | Política por defecto ante duda o invalidación: **bloquear y escalar a humano**. Nunca continuar con error. | "Error", "excepción" |

---

## 3. Actores y autoridad sobre hechos

| Actor | Tipo | Hechos sobre los que tiene autoridad |
|---|---|---|
| Paciente/Socio | Humano | Otorgar o rechazar el consentimiento de retiro |
| EnfermeroRetirador | Humano | Detección, lectura de MedicationIdentity, custodia física hasta base |
| FarmacéuticoClasificador | Humano (DT regulatorio) | ClasificacionResiduo · firma del Manifiesto · cierre de cadena |
| MédicoTratante | Humano | Reemplazo terapéutico si el medicamento retirado estaba activo |
| TransportistaHabilitado | Organización | Recepción y transporte bajo manifiesto |
| OperadorDisposicion | Organización | Emisión del certificado de disposición final |
| AuditorRegulatorio | Humano (solo lectura) | Consulta completa de la cadena |
| SistemaGentleCare | Software | Detectar, leer, notificar, emitir receipts. **Sin autoridad sobre hechos clínicos ni regulatorios.** |

**Invariante de autoridad:** ningún hecho regulatorio (clasificación, manifiesto, disposición) puede ser autorado por el SistemaGentleCare.

---

## 4. Ciclo de vida de `RetiroFarmaceutico` (máquina de estados)

```
Detectado
  → ConsentimientoGestionado {Obtenido | Rechazado}
  → Clasificado
  → EnTransitoABase
  → EnBaseConsolidado
  → ManifiestoEmitido
  → EnTransitoTransportista
  → RecibidoEnPlanta
  → DisposicionFinalCertificada
  → Cerrado
```

**Ramas terminales:** `RechazadoPorPaciente` (cierre con registro) · `BloqueadoFailClose` (requiere resolución humana).

| Desde | Hacia | Comando | Guard |
|---|---|---|---|
| Detectado | ConsentimientoGestionado | SolicitarConsentimiento | MedicationIdentity completa (I6) |
| ConsentimientoGestionado | Clasificado | ClasificarResiduo | ConsentimientoObtenido (I1) |
| Clasificado | EnTransitoABase | ConsolidarEnBase | Compartimento inverso sellado |
| EnBaseConsolidado | ManifiestoEmitido | EmitirManifiesto | Farmacéutico responsable (I3) |
| ManifiestoEmitido | EnTransitoTransportista | TransferirATransportista | Habilitación OPDS vigente al momento (I4) + firma dual (I2) |
| EnTransitoTransportista | RecibidoEnPlanta | RecibirEnPlanta | Firma dual (I2) + sin saltos (I5) |
| RecibidoEnPlanta | DisposicionFinalCertificada | CertificarDisposicionFinal | Operador habilitado vigente (I4) |
| DisposicionFinalCertificada | Cerrado | CerrarEvento | Cadena completa sin saltos (I5) |
| Cualquier estado | BloqueadoFailClose | BloquearFailClose | Automático ante invalidación (I9) |

---

## 5. Comandos

| Comando | Actor con autoridad |
|---|---|
| DetectarMedicamento | EnfermeroRetirador |
| LeerMedicationIdentity | Sistema (en nombre de EnfermeroRetirador) |
| SolicitarConsentimiento / RegistrarConsentimiento / RegistrarRechazo | Paciente/Socio |
| ClasificarResiduo | FarmacéuticoClasificador |
| ConsolidarEnBase | EnfermeroRetirador |
| EmitirManifiesto | FarmacéuticoClasificador |
| TransferirATransportista | FarmacéuticoClasificador → TransportistaHabilitado |
| RecibirEnPlanta | OperadorDisposicion |
| CertificarDisposicionFinal | OperadorDisposicion |
| CerrarEvento | FarmacéuticoClasificador |
| BloquearFailClose | SistemaGentleCare (automático) |
| ResolverBloqueo | FarmacéuticoClasificador o Director Médico |

---

## 6. Eventos de dominio (payload mínimo)

| Evento | Payload mínimo |
|---|---|
| `MedicamentoVencidoDetectado` | eventoId · MedicationIdentity[] · domicilioRef · lector · timestamp |
| `ConsentimientoObtenido` | eventoId · pacienteRef · receiptHash |
| `ConsentimientoRechazado` | eventoId · motivo |
| `CadenaDeFrioRota` | unidadRef · rangoObservado · rangoEsperado (+2 °C a +8 °C) |
| `RecallANMATActivado` | gtin · lote · alertaRef |
| `ResiduoClasificado` | eventoId · clase · farmacéuticoId |
| `ManifiestoGenerado` | manifestId · triplicadoHash |
| `TransferenciaATransportista` | nodoId · firmas[2] · lugar · tiempo |
| `RecepcionEnPlanta` | nodoId · operadorId · firmas[2] |
| `DisposicionFinalCertificada` | certificadoId · método · hash |
| `FailCloseActivado` | motivo · regla · estadoBloqueado |
| `EventoCerrado` | eventoId · cadenaCompleta: bool |

---

## 7. Invariantes

- **I1** — No hay retiro sin ConsentimientoObtenido. Excepción recall ANMAT: régimen específico `[VALIDAR-LEGAL]` (ver R1).
- **I2** — Todo NodoDeTransferencia tiene exactamente dos firmas: quien entrega y quien recibe.
- **I3** — Todo Manifiesto exige FarmacéuticoClasificador como responsable técnico.
- **I4** — Transportista y Operador con habilitación **vigente al momento del nodo** (validación temporal, no "actual").
- **I5** — CadenaDeCustodia sin saltos: el nodo N+1 comienza donde terminó el N.
- **I6** — MedicationIdentity completa (4 componentes) o Fail-Close.*[HIPÓTESIS — pendiente de validación farmacéutica; no es invariante exigible hasta el Gate E1]*
- **I7** — Una unidad retirada **nunca** vuelve al circuito comercial ni al stock válido.
- **I8** — Todo cambio de estado emite exactamente un evento de dominio.
- **I9** — Fail-Close por defecto: ante duda o invalidación, bloquear y escalar.

---

## 8. Failure paths y políticas

| Failure path | Detección | Política | Responsable de resolución |
|---|---|---|---|
| Paciente rechaza el retiro | RegistrarRechazo | Cierre con registro + notificación al médico si el medicamento estaba activo | EnfermeroRetirador |
| Identidad no legible o no reconocida en ANMAT | Lectura fallida | Fail-Close → clasificación manual | FarmacéuticoClasificador |
| Cadena de frío rota en compartimento inverso | Telemetría IoT | Alerta + bloqueo de transferencia hasta validación | FarmacéuticoClasificador |
| Transportista no arriba en 72 h | Timeout de nodo | Escalación a Director Médico + transportista alternativo del registro vigente | Director Médico |
| Operador no certifica en 30 días | Timeout de certificado | Alerta regulatoria interna; evaluar reporte a autoridad `[VALIDAR-LEGAL]` (R4) | FarmacéuticoClasificador |
| Salto en cadena de custodia (nodo faltante) | Verificación I5 | Fail-Close + auditoría interna inmediata | Director Médico |
| Recall con paciente no localizable | Protocolo de intentos | 3 intentos documentados por canales alternos `[VALIDAR-LEGAL]` | EnfermeroRetirador + MédicoTratante |

---

## 9. Sistemas externos y fronteras

| Sistema | Propiedad | Interacción | Frontera |
|---|---|---|---|
| ANMAT — Trazabilidad (Disp. 3683/2011) | Externa | Consulta de estado GTIN+serie y recalls | Solo lectura |
| OPDS — Registro transportistas/plantas PBA | Externa | Validez de habilitaciones | Solo lectura + cache con refresh |
| Google Cloud Healthcare API (FHIR) | Alianza | Publicación de eventos clínicos relevantes al registro del paciente | Escritura con consentimiento |
| Notificación a MédicoTratante | Nurse Canning 24 | Alertas de retiro de medicamento activo | Canal interno |
| Telemetría IoT cadena de frío | Nurse Canning 24 | Temperatura del compartimento inverso | Solo lectura |

**Nota de frontera:** ningún sistema externo posee `RetiroFarmaceutico`; el aggregate vive en Gentle Care.

---

## 10. Semántica del ReceiptOperacional — decisión abierta

### 10.1 Qué demuestra SHA-256 (y qué no)

| Propiedad | ¿SHA-256 solo la demuestra? | Qué se necesitaría además |
|---|---|---|
| Integridad del contenido | ✅ Sí | — |
| Identidad del autor | ❌ No | Clave ligada a identidad (certificado/PKI) |
| Autenticidad / origen | ❌ No | Firma digital con clave privada del actor |
| Timestamp confiable | ❌ No | TSA (RFC 3161) o log confiable |
| No repudio | ❌ No | Firma + custodia de claves + política de certificación |
| Inmutabilidad en el tiempo | ❌ No | Store append-only / anclaje (transparency log) |
| Cumplimiento legal | ❌ No | Marco normativo + firmas de partes competentes + retención |

### 10.2 Hechos a certificar y autoridad

| Hecho | Autoridad (ver §3) | Evidencia ligada |
|---|---|---|
| Consentimiento de retiro | Paciente/Socio | Texto informado + firma del paciente + timestamp |
| Clasificación del residuo | FarmacéuticoClasificador | Clase + criterio + identidad del farmacéutico |
| Cada transferencia de custodia | Ambos actores del nodo | Firma dual + lugar + tiempo + lectura de unidades |
| Disposición final | OperadorDisposicion | Certificado con método y hash de manifiesto |

### 10.3 Decisión diferida al Engineering Gate

Esquema de firma, timestamping y anclaje **no se deciden en discovery**. Hasta el Gate: sin nombre de paquete, sin implementación.
**Distinción explícita:** `DevelopmentReceipt` (RDD de Gentle-AI, evidencia del proceso de desarrollo) ≠ `ReceiptOperacional` (Gentle Care, evidencia de eventos clínicos). Contextos y lifecycles distintos: **no compartir paquete por asumir**.

---

## 11. Preguntas regulatorias abiertas (asesor legal)

- **R1** — En recall ANMAT, ¿procede el retiro sin consentimiento o con notificación reforzada?
- **R2** — Retención mínima por tipo documental (manifiesto, certificado, consentimiento): se propone 10 años, confirmar.
- **R3** — Validez de firma digital del paciente para consentimiento (Ley 25.326, Ley 26.529, normativa PBA).
- **R4** — Obligación de reporte a OPDS/autoridad ante operador sin certificado en plazo.
- **R5** — Resolución vigente de clasificación de residuos peligrosos/especiales PBA (número y vigencia).
- **R6** — Carácter de "generador" de residuos del JV (Clínica Canning Health): confirmar titularidad.

---

## 12. Candidatos de agregados / entidades / value objects (hipótesis, no decisiones)

- `RetiroFarmaceutico` (aggregate root) conteniendo `CadenaDeCustodia` como colección de `NodoDeTransferencia` (entities).
  - *Hipótesis alternativa:* CadenaDeCustodia como aggregate propio. **Criterio de decisión para el Gate:** ¿comparten lifecycle y consistencia transaccional en el mismo boundary?
- `MedicationIdentity`, `ClasificacionResiduo`, `FirmaDual` → value objects.
- `OperadorHabilitado` → aggregate de datos maestros (invariante: vigencia de habilitación).
- `Manifiesto` → entity dentro de RetiroFarmaceutico **vs** entity de un contexto regulatorio separado: **abierto**.
- **No decidido:** nombres de paquetes, persistencia, transporte, formatos de firma.

---

## 13. Fuera de alcance MVP-0

| Flujo excluido | Razón |
|---|---|
| Residuos patogénicos | Otro régimen (Ley 24.049 + Dec. 403/97 PBA) y otro lifecycle → v2.0 |
| Estupefacientes y psicotrópicos | Ley 17.565 Art. 12 + Sedronar |
| Citostáticos oncológicos | Bioseguridad nivel 2 |
| Residuos radiactivos | Competencia ARN |

---

## 14. Engineering Gate — criterios de aceptación

- [ ] Revisión de Farmacéutico Institucional (glosario §2, clasificación §5)
- [ ] Validación de asesor legal ambiental-sanitario (§11, R1-R6)
- [ ] Revisión técnica (ciclo de vida §4, invariantes §7, semántica de receipts §10)
- [ ] Decisiones de §10 y §12 registradas como ADR
- [ ] Solo después de esto: derivar Issues de implementación

| Rol | Nombre | Firma | Fecha |
|---|---|---|---|
| Owner de producto/dominio | Barto912 | | |
| Farmacéutico Institucional | | | |
| Asesor legal ambiental-sanitario | | | |
| Revisión técnica | | | |

---

## Anexo A — Trazabilidad con PROTO-006-LIF v1.1

| Sección de este documento | Sección del protocolo |
|---|---|
| §1 Alcance | §2.1-2.5 |
| §2 Lenguaje ubicuo | §4 Definiciones |
| §3 Actores | §5 Matriz RACI |
| §4 Ciclo de vida | §6 Flujos A-C |
| §7 Invariantes | §8 Triggers Fail-Close |
| §8 Failure paths | §8 + §12 Gobernanza |
| §9 Sistemas externos | §2.3 + §9 |
| §10 Receipts | §9.1 Receipt SHA-256 |
