# Aggregate Boundaries — MVP-0 Logística Inversa Farmacéutica

**Bounded context:** `logistica-inversa-farmaceutica`
**Versión:** 0.2 (propuesta ajustada según revisión de Deivis — PASS WITH CHANGES)
**Cambios v0.2:** (a) I8 reformulada sin cardinalidad 1:1 command/event; (b) ownership neutral `ResponsableRegulatorio`; (c) inputs diferidos H13/H14 registrados.
**Fecha:** 7 de octubre de 2026 (v0.1) · 8 de octubre de 2026 (v0.2)
**Autor:** Barto912 (owner de producto/dominio)
**Base:** `logistica-inversa.md` v0.2 · `event-storming-mvp0.md` (hotspots H1-H12) · revisión E0 (Deivis) · revisión de Aggregate Boundaries (Deivis): PASS WITH CHANGES
**Estado de madurez:** architecture discovery — sin código de producción hasta el Architecture Gate E1.

---

## 0. Convenciones de modelado (resuelve H11)

| Elemento | Regla | Ejemplo |
|---|---|---|
| Command | Verbo imperativo | `IniciarRetiro`, `BloquearFailClose`, `ResolverBloqueo` |
| Domain event | Pasado, hecho ocurrido | `RetiroIniciado`, `FailCloseActivado`, `BloqueoResuelto` |
| State | Nodo del lifecycle del aggregate | `Bloqueado`, `Cerrado`, `Cancelado` |

**Renames acordados:** el state `BloqueadoFailClose` pasa a llamarse `Bloqueado`. `BloqueadoFailClose` **no es un evento** y no debe aparecer como tal en ningún documento del contexto. El evento es `FailCloseActivado`; el command es `BloquearFailClose`.

**I8 reformulada (revisión de Deivis):** un command puede producir 0..N Domain Events dentro de una misma decisión/transacción. I8 exige que toda transición relevante quede representada por uno o más eventos suficientes para reconstrucción y auditoría, sin imponer cardinalidad 1:1 entre command y event.

---

## 1. Inputs de la revisión E0 y dónde se resuelven

| Hotspot | Dónde se resuelve | Estado |
|---|---|---|
| H9 Trigger del aggregate | §2.1 | Resuelto en esta propuesta |
| H10 Semántica de cierre | §2.2 | Resuelto en esta propuesta |
| H11 Command / event / state | §0 | Resuelto en esta propuesta |
| H1 CadenaDeCustodia | §3.1 | Resuelto en esta propuesta |
| H2 Manifiesto | §3.2 | Resuelto en esta propuesta |
| H12 Frontera FHIR / Cloud | Context Map | Diferido (patrón Integration Contract/ACL aceptado) |

---

## 2. Aggregate root: `RetiroFarmaceutico`

### 2.1 Trigger y birth event (H9)

`MedicamentoVencidoDetectado` **no** es el evento inicial general: es uno de varios triggers posibles. El aggregate nace con:

- **Command:** `IniciarRetiro(causa, MedicationIdentity[])`
- **Birth event:** `RetiroIniciado`, con `causa ∈ {VENCIMIENTO, DISCONTINUACION, CADENA_FRIO_COMPROMETIDA, RECALL_ANMAT}`

Los eventos de detección (`MedicamentoVencidoDetectado`, `CadenaDeFrioRota`, `RecallANMATActivado`, `MedicamentoDiscontinuadoDetectado`) quedan como **eventos trigger upstream**: políticas reactivas los convierten en el command `IniciarRetiro`. Esto preserva I8 reformulada (toda transición relevante representada por eventos suficientes para reconstrucción y auditoría) y evita que el aggregate dependa de una única causa de origen.

### 2.2 Lifecycle y semántica de cierre (H10)

States: `Iniciado → ConsentimientoGestionado → Clasificado → Consolidado → ManifiestoEmitido → EnTransito → RecibidoEnPlanta → DisposicionCertificada`, más el state no terminal `Bloqueado`.

**Tres resultados terminales distintos, con eventos distintos:**

| Terminal | Event | Condición | Semántica |
|---|---|---|---|
| Cierre exitoso | `RetiroCerrado` | Cadena de custodia completa sin saltos (I5) | Disposición final certificada |
| Cancelación | `RetiroCancelado` | Sin retiro físico: `ConsentimientoRechazado` o cancelación antes de consolidación | El retiro no ocurrió; queda registro auditable |
| Cierre por excepción | `RetiroCerradoPorExcepcion` | `Bloqueado` resuelto por autoridad humana (FarmacéuticoClasificador o Director Médico) con cadena incompleta documentada | El retiro ocurrió parcialmente; la excepción queda firmada y auditada |

**Razón:** mezclar rechazo y cierre exitoso en un mismo evento destruiría la capacidad de auditoría (no se podría distinguir operacionalmente un retiro completado de uno abortado) y vaciaría de contenido a I5.

### 2.3 Contenido del aggregate

- **Entities:** `NodoDeTransferencia` (colección `CadenaDeCustodia`, §3.1) · `Manifiesto` (relación 1:1, §3.2)
- **Value objects:** `MedicationIdentity` (GTIN + serie + lote + vencimiento) · `ClasificacionResiduo` · `FirmaDual` · `Consentimiento`
- **Invariantes propias:** I1-I9 del discovery v0.2 (I6 como hipótesis pendiente de validación farmacéutica)

---

## 3. Decisiones de boundary

### 3.1 H1 — `CadenaDeCustodia`: entities dentro del aggregate

**Decisión:** colección de entities `NodoDeTransferencia` dentro de `RetiroFarmaceutico`. **No** es aggregate propio.

Justificación contra los cuatro criterios de revisión:
- *Invariantes:* I2 (firma dual) e I5 (sin saltos) se validan **sobre el conjunto de nodos**, no sobre un nodo aislado: exigen visibilidad del conjunto en la misma transacción.
- *Consistencia transaccional:* agregar un nodo y verificar continuidad con el nodo anterior es una sola transacción; separar el aggregate obligaría a sagas distribuidas sin beneficio.
- *Ownership:* un nodo no existe sin su retiro; no tiene identidad independiente ni escritor externo.
- *Failure boundary:* un nodo inválido (falta firma, habilitación vencida) dispara `FailCloseActivado` **sobre el mismo aggregate**, sin coordinación entre boundaries.

### 3.2 H2 — `Manifiesto`: entity 1:1 dentro del aggregate (MVP-0)

**Decisión:** entity dentro de `RetiroFarmaceutico` mientras la relación sea 1:1.

Justificación:
- *Consistencia transaccional:* la emisión del manifiesto es **guard** de `TransferirATransportista`; vive en la misma transacción que el nodo de transferencia que lo referencia.
- *Ownership:* en MVP-0 el manifiesto es emitido y custodiado por este contexto (firmado por FarmacéuticoClasificador, I3).
- *Failure boundary:* manifiesto inválido = transferencia bloqueada dentro del mismo boundary.

**Criterio de reversión (explícito):** si en v2 aparecen manifiestos consolidados (N retiros por ruta/manifiesto) o el transportista requiere lifecycle propio del documento, se promueve `Manifiesto` a aggregate root en un bounded context `regulatorio-transporte`. Esa promoción será un ADR, no un refactor silencioso.

### 3.3 Aggregate separado: `OperadorHabilitado` (master data)

**Decisión:** aggregate root propio, fuera de `RetiroFarmaceutico`.

Justificación:
- *Lifecycle:* la vigencia de habilitaciones OPDS evoluciona independientemente de cualquier retiro.
- *Ownership:* es master data regulatoria con otro escritor (sincronización desde registro externo, ver ADR-003).
- *Consistencia:* `RetiroFarmaceutico` lo consume por referencia (id + cache con refresh) en el guard I4; **nunca** en la misma transacción.
- *Failure boundary:* habilitación vencida detectada en un transfer → `FailCloseActivado` sobre `RetiroFarmaceutico`; `OperadorHabilitado` no se modifica.

**Pendiente (H13):** determinar en el Context Map si `OperadorHabilitado` pertenece realmente a este bounded context o es una proyección/reference model proveniente de un Regulatory Registry Context.

---

## 4. Consistencia transaccional y failure boundaries

| Relación | Modelo |
|---|---|
| Dentro de `RetiroFarmaceutico` | Consistencia fuerte: un command = una transacción = 0..N events suficientes para reconstrucción y auditoría (I8 reformulada) |
| Entre aggregates (`RetiroFarmaceutico` ↔ `OperadorHabilitado`) | Consistencia eventual por referencia + cache con refresh; validación temporal en el guard (I4) |
| Con sistemas externos (ANMAT, OPDS, FHIR, IoT) | **Fuera de toda transacción**: Integration Contract/ACL + timeouts + Fail-Close (H12, Context Map) |

**Matriz de failure boundary:**

| Falla | Aggregate afectado | Política |
|---|---|---|
| Firma dual incompleta en nodo | RetiroFarmaceutico | `FailCloseActivado` (I2) |
| Salto de cadena detectado | RetiroFarmaceutico | `FailCloseActivado` (I5) |
| Habilitación vencida en transfer | RetiroFarmaceutico | `FailCloseActivado` (I4) |
| Refresh de cache OPDS caído | RetiroFarmaceutico | Bloqueo preventivo de transfers hasta refresh (P6) |
| Timeout 72 h sin transferencia | RetiroFarmaceutico | Escalada a Director Médico (P4) |
| Datos maestros de operador corruptos | OperadorHabilitado | Cuarentena del registro + alerta; no afecta retiros en curso |

---

## 5. Ownership y autoridad por aggregate

| Aggregate | Owner humano | Hechos que autoriza |
|---|---|---|
| `RetiroFarmaceutico` | FarmacéuticoClasificador (DT) + actores del flujo (discovery §3) | Clasificación, manifiesto, transferencias, cierre y excepciones |
| `OperadorHabilitado` | Director Médico / `ResponsableRegulatorio` | Alta, baja y vigencia de habilitaciones sincronizadas |

**Nota de neutralidad (revisión de Deivis):** mientras el JV siga siendo PROPOSED/TARGET, no forma parte del modelo de dominio; los roles se expresan de forma neutral (`ResponsableRegulatorio`).

El SistemaGentleCare **no es owner de ningún aggregate**: emite commands asistenciales y receipts, sin autoridad sobre hechos clínicos o regulatorios (invariante de autoridad, discovery §3).

---

## 6. Alineación pendiente (discovery v0.3 y event storming v0.2, post-aprobación)

- Events nuevos: `RetiroIniciado`, `RetiroCancelado`, `RetiroCerradoPorExcepcion`.
- Renames: `EventoCerrado` → `RetiroCerrado` · `CerrarEvento` → `CerrarRetiro` · state `BloqueadoFailClose` → `Bloqueado`.
- Command nuevo: `IniciarRetiro`.
- I6 **no se toca**: sigue como hipótesis hasta validación farmacéutica.

---

## 7. Checklist de revisión (criterios de Deivis)

- [x] *Invariantes:* cada invariante I1-I9 mapeada a un aggregate único que la exige (§2.3, §3.3); I8 reformulada sin cardinalidad 1:1 (§0, §4)
- [x] *Consistencia transaccional:* ningún invariante cruzando boundaries sin justificación (§4)
- [x] *Ownership:* todo hecho con un único owner humano identificado, con roles neutrales (§5)
- [ ] *Failure boundaries:* toda falla con aggregate afectado y política explícita (§4) — pendiente H14: distinguir reparación vs. aceptación documentada de excepción en `ResolverBloqueo`

---

## 8. Próximos pasos del recorrido

1. Context Map con H12 y H13: `Logística Inversa → Integration Contract/ACL → FHIR → infraestructura`; determinar si `OperadorHabilitado` es aggregate de este contexto o proyección de un Regulatory Registry Context.
2. Resolver H14 antes de cerrar ADR-001: `ResolverBloqueo` ¿repara la condición que violó I5 o acepta/documenta una excepción todavía existente? La autoridad humana puede cerrar operacionalmente una excepción, pero no convierte retroactivamente una cadena inválida en válida.
3. ADR-001 (aggregate boundaries), ADR-002 (operational evidence), ADR-003 (registros externos y offline).
4. Architecture Gate E1.
