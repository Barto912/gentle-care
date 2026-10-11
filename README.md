# NurseHealth.ai

**Enfermería domiciliaria de precisión con IA en el borde.**

[![CI](https://github.com/Barto912/gentle-care/actions/workflows/rdd-ci.yml/badge.svg)](https://github.com/Barto912/gentle-care/actions)

> 🚧 **Estado del proyecto: INVESTIGACIÓN Y DISEÑO.**
> Este repositorio no es un producto en producción.
> El código y los documentos aquí presentes son experimentales
> y su uso en entornos clínicos reales no está autorizado ni validado.

## El problema

Miles de pacientes crónicos en las Américas reciben atención
domiciliaria sin trazabilidad real. Ese trabajo hoy es invisible:
planillas de papel, medicamentos vencidos en el hogar, cadenas de
custodia rotas y decisiones clínicas sin evidencia auditable.

**Piloto inicial:** 57 countries y barrios cerrados en Canning,
Argentina.
**Visión:** escalamiento regional a Chile, Uruguay, Colombia,
México y Brasil con el mismo protocolo.

## La solución

Plataforma de enfermería de precisión extra-hospitalaria:

- **IA en el borde:** identificación de medicamentos en el
  dispositivo del enfermero, sin enviar datos del paciente a la nube.
- **Cadena de custodia criptográfica:** cada evento sellado con
  SHA-256 en almacenamiento append-only.
- **Logística inversa farmacéutica:** retiro domiciliario de
  medicamentos vencidos con trazabilidad hasta disposición final.
- **Fail-close por diseño:** ante la duda, el sistema bloquea y
  escala a autoridad humana. Nunca improvisa.

**Lema:** el enfermero decide, el sistema prueba.

## Arquitectura

```text
Pixel 11 Pro (edge AI, TensorFlow Lite)
        |
Cloud Run (Go, API REST, event sourcing)
        |
Cloud SQL (PostgreSQL, event store)
        |
Google Healthcare API (FHIR R4, detras de ACL)
```

## Stack técnico

| Capa | Tecnología |
| --- | --- |
| Backend | Go en Cloud Run |
| Frontend | React Native con Expo |
| Base de datos | PostgreSQL en Cloud SQL |
| IA en el borde | TensorFlow Lite en Pixel 11 Pro |
| Interoperabilidad | FHIR R4 vía Google Healthcare API |
| Infraestructura | Google Cloud Platform |

## Estado actual

Diseño completado y documentado:

- Event Storming del MVP-0 con 14 hotspots resueltos.
- Aggregate Boundaries v0.3 (reparación vs aceptación documentada).
- Context Map con ACL para FHIR y Reference Model regulatorio.
- Integración continua verde en cada commit.

En curso: implementación del MVP con créditos de
Google for Startups Cloud Program.

## Gobernanza ética

Principios no negociables del manifiesto del proyecto:

1. Dignidad sobre datos.
2. Autonomía sobre automatización.
3. Transparencia sobre caja negra.
4. Cuidado sobre eficiencia.
5. Presencia sobre conexión.
6. Prevención sobre reacción.
7. Formación sobre adoctrinamiento.

Los datos de pacientes solo se procesan seudonimizados.
Ninguna decisión clínica es automatizada.

## Equipo

- **Barto** — Founder y product owner.
- **Deivis** — Technical advisor (DDD, arquitectura event-driven).
- **Gabriel** — Chief Medical Officer.
- **Leonor** — Strategic advisor.

## Licencia

Propietario. Todos los derechos reservados.
