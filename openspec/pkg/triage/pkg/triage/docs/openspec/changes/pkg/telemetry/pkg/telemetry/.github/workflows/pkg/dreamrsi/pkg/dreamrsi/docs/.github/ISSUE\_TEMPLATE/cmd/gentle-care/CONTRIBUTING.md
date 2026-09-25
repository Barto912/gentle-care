```
# 🤝 Guía de Contribución a Gentle Care

Agradecemos tu interés en colaborar con **Gentle Care**. Para garantizar la precisión clínica, la privacidad de datos (**PHI/PII**) y la estabilidad técnica, todas las contribuciones se rigen por la arquitectura de **Gentle-AI v3.7.0**.

---

## 🚦 Vías de Desarrollo

### 1. Organic-Driven Development (ODD) - Vía Ágil
* **Aplica para:** Correcciones menores de interfaz, ajustes de logs, refactorización ligera o documentación.
* **Proceso:** Abre un *Issue* o *Pull Request* directo. La memoria del agente en **Gentle Engram™** preservará el contexto acumulado entre revisiones.

### 2. Spec-Driven Development (SDD) - Vía Clínica Formal
* **Aplica para:** Modificaciones en algoritmos de Triage, lógica médica, taxonomía **SNOMED CT**, telemetría IoT o ciberseguridad.
* **Proceso:** Requiere crear previamente un documento **OpenSpec** en `openspec/changes/` detallando la propuesta y los criterios de prueba antes de escribir código.

---

## 🛡️ Aduana de Validación Criptográfica (RDD &amp; Strict TDD)
1. **Tests Obligatorios:** Todo código debe incluir pruebas unitarias en Go (`*_test.go`).
2. **Prueba 4R:** Los revisores auditarán la propuesta bajo los criterios de **Risk** (seguridad), **Readability** (legibilidad), **Reliability** (confiabilidad determinista) y **Resilience** (resiliencia de integración).
3. **Receipt SHA-256:** Ningún código se despliega sin un recibo inmutable verificado bajo la política **Fail-Close**.

*Hecho con Gentle-AI.*

```
