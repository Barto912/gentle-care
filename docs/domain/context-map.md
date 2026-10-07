# Context Map — MVP-0 Logística Inversa Farmacéutica

**Bounded context:** `logistica-inversa-farmaceutica`
**Versión:** 0.1 (propuesta — pendiente de revisión técnica)
**Fecha:** 8 de octubre de 2026
**Autor:** Barto912 (owner de producto/dominio)
**Base:** `logistica-inversa.md` v0.2 · `event-storming-mvp0.md` · `aggregate-boundaries.md` v0.2 · revisión E0 y revisión de Aggregate Boundaries (Deivis)
**Estado de madurez:** architecture discovery — sin código de producción hasta el Architecture Gate E1.

---

## 0. Resumen de decisiones de Context Map

| Hotspot | Decisión | Patrón DDD |
|---|---|---|
| **H12** | `Google Cloud Healthcare API` y FHIR quedan **fuera** del bounded context. Se accede vía Anti-Corruption Layer (ACL). | ACL + Published Language |
| **H13** | `OperadorHabilitado` **no** es aggregate de este contexto. Es un Reference Model (proyección) del `Regulatory Registry Context`. | Customer-Supplier + Read Model |

---

## 1. Mapa de contextos (big picture)

```
┌────────────────────────────────┐        ┌────────────────────────────────┐
│  Regulatory Registry Context   │        │     IoT Telemetry Context      │
│  (OPDS / ANMAT)                │        │     (cadena de frío)           │
└───────────────┬────────────────┘        └───────────────┬────────────────┘
                │ Customer-Supplier                       │ Conformist
                │ (proyección / cache)                    │ (eventos de telemetría)
                ▼                                         ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                    logistica-inversa-farmaceutica                          │
│                    (RetiroFarmaceutico = aggregate root)                   │
└───────────────┬────────────────────────────────────────┬───────────────────┘
                │ ACL + Published Language               │ Open Host Service
                │ (eventos de dominio → FHIR R4)         │ (alertas / notificaciones)
                ▼                                        ▼
┌────────────────────────────────┐        ┌────────────────────────────────┐
│   Clinical EHR / FHIR Context  │        │     Canal Médico Tratante      │
│   (Google Cloud Healthcare API)│        │                                │
└────────────────────────────────┘        └────────────────────────────────┘
```

---

## 2. Resolución de H12: frontera FHIR / Cloud (ACL)

**Problema:** el discovery original mencionaba `Google Cloud Healthcare API` como sistema externo, lo que rischiaba acoplar el modelo de dominio al proveedor de infraestructura.

**Decisión (patrón ACL — Anti-Corruption Layer):** el bounded context `logistica-inversa-farmaceutica` **no conoce** FHIR, HL7 ni Google Cloud. Solo conoce sus propios eventos de dominio (`RetiroCerrado`, `RetiroCancelado`, etc.).

**Flujo de integración:**
1. El aggregate `RetiroFarmaceutico` emite un Domain Event (ej. `RetiroCerrado`).
2. Un **Integration Contract / ACL** interno escucha el evento.
3. El ACL traduce el evento de dominio al estándar FHIR R4 (ej. `MedicationStatement` u `Observation`).
4. El ACL invoca la infraestructura (Google Cloud Healthcare API) para persistir el recurso FHIR.

**Beneficio:** si mañana cambiamos de Google Cloud a otro proveedor, o de FHIR a un legacy HL7 v2, el cambio ocurre **exclusivamente en el ACL**. El modelo de dominio y los aggregates permanecen intactos.

---

## 3. Resolución de H13: `OperadorHabilitado` (Reference Model)

**Problema:** en `aggregate-boundaries.md` v0.1, `OperadorHabilitado` estaba modelado como Aggregate Root separado dentro de este bounded context. La revisión cuestionó su ownership al ser master data sincronizada desde un source of truth externo.

**Decisión (patrón Customer-Supplier + Read Model):** `OperadorHabilitado` **no pertenece** a `logistica-inversa-farmaceutica`. La fuente de verdad (Source of Truth) es el registro público de OPDS / ANMAT.

**Modelado correcto:**
- Existe un contexto upstream llamado **`Regulatory Registry Context`** (servicio de sincronización con los registros públicos).
- Nuestro contexto consume esos datos y mantiene una **proyección local (Read Model / cache)** llamada `OperadorHabilitadoCache`.
- **Invariante I4 (habilitación vigente):** se valida contra la proyección local. Si la cache está caída o vencida, se aplica la política P6 (bloqueo preventivo), pero **nunca** se escribe ni modifica el registro upstream desde este contexto.

**Impacto en Aggregate Boundaries:** `OperadorHabilitado` deja de ser Aggregate Root de este contexto y pasa a ser un **Reference Model** inyectado en el guard de `TransferirATransportista`.

---

## 4. Contratos de integración (Integration Contracts)

### 4.1 Regulatory Registry (upstream)
- **Fuente:** OPDS / ANMAT.
- **Patrón:** Customer-Supplier (somos clientes de su API pública).
- **Contrato:** sincronización asíncrona (batch diario o webhook).
- **Fallback:** si el registro upstream cae, se usa la última proyección válida (TTL 24 h). Si el TTL expira, Fail-Close: no se permiten transferencias.

### 4.2 IoT Telemetry (upstream)
- **Fuente:** sensores de caja térmica.
- **Patrón:** Conformist (consumimos los eventos tal como los emite el hardware).
- **Contrato:** stream de eventos (`TemperatureReading`, `SealBroken`).
- **Política:** si el stream se interrumpe por más de 15 min durante un traslado, `CadenaDeFrioRota` (Fail-Close preventivo).

### 4.3 Clinical EHR / FHIR (downstream)
- **Destino:** Historia Clínica Electrónica (vía Google Cloud Healthcare API).
- **Patrón:** Published Language + ACL.
- **Contrato:** nuestro contexto publica eventos de dominio (`RetiroCerrado`, `RetiroCancelado`). El ACL traduce a FHIR R4 (`MedicationStatement` con status `completed` o `entered-in-error`).
- **Consistencia:** eventual, con **Outbox Pattern**: si la API de FHIR cae, el evento se encola y se reintenta. El retiro físico **no se bloquea** por una caída del EHR (prioridad operativa).

### 4.4 Canal Médico Tratante (downstream)
- **Destino:** médico del paciente (notificaciones P1 y P4).
- **Patrón:** Open Host Service / Notification.
- **Contrato:** API REST simple o gateway de notificaciones. Solo lectura: el médico no responde vía API, actúa en el mundo real o en su propio EHR.

---

## 5. Matriz de dependencias y failure boundaries

| Contexto vecino | Dependencia crítica | Comportamiento ante fallo (failure boundary) |
|---|---|---|
| **Regulatory Registry** | Alta (I4) | Bloqueo de transferencias si el TTL de cache expira. No afecta retiros ya en curso. |
| **IoT Telemetry** | Alta (cadena de frío) | Fail-Close automático si se pierde telemetría en tránsito. |
| **Clinical EHR (FHIR)** | Baja (operativa) | Outbox Pattern: el retiro se cierra operativamente; la sincronización clínica se reintenta en background. |
| **Canal Médico** | Baja (informativa) | Reintento de notificación. No bloquea el flujo físico. |

---

## 6. Alineación pendiente (aggregate-boundaries.md v0.3, post-aprobación)

- Mover `OperadorHabilitado` de la sección "Aggregate separado" a "Reference Models / proyecciones".
- Explicitar el Outbox Pattern para la integración con FHIR (ACL).

---

## 7. Próximos pasos del recorrido

1. Revisión de este Context Map (Deivis).
2. Resolución de **H14** (semántica de `ResolverBloqueo`: reparación vs. aceptación documentada de excepción) antes de cerrar ADR-001.
3. ADR-001 (aggregate boundaries), ADR-002 (operational evidence), ADR-003 (registros externos y offline).
4. Architecture Gate E1.
