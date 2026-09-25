```
# OpenSpec: Dream-RSI Triage Replay &amp; Monotonic Safety Gate

- **ID:** `change-20260925-dream-rsi-triage`
- **Target:** Gentle Care Engine v3.7.0 + Google Pixel 11 Pro
- **Status:** Proposed

## Requisitos de Seguridad Médica
1. **R-DREAM-01 (Zero-Execution Replay):** El simulador DEBE evaluar variaciones de la política de triage reordenando el árbol histórico de **Gentle Engram** a costo cero de inferencia de LLM.
2. **R-DREAM-02 (Never-Worse Gate):** Ninguna actualización del orquestador de triage podrá desplegarse si su puntaje en el simulador histórico es inferior al de la versión vigente (`Score(pi_t+1) &gt;= Score(pi_t)`).
3. **R-DREAM-03 (Edge Optimization):** Reducción obligatoria del 30% al 50% en consumo de batería e inferencia remota en el dispositivo Pixel 11 Pro para eventos de baja complejidad (**Código Celeste**).

```
