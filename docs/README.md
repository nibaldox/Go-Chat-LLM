# TUI para Ollama en Go

# Chat LLM TUI

Interfaz de texto (TUI) en Go para conversar de forma interactiva con modelos locales servidos por **Ollama** (`http://localhost:11434`). Inspirada en la experiencia de ChatGPT, pero totalmente offline y optimizada para terminal.

---

## Tabla de contenido
1. [Requisitos](#requisitos)
2. [Instalación rápida](#instalación-rápida)
3. [Uso](#uso)
4. [Atajos de teclado](#atajos-de-teclado)
5. [Arquitectura](#arquitectura)
6. [Desarrollo y mantenimiento](#desarrollo-y-mantenimiento)
7. [Hoja de ruta](#hoja-de-ruta)
8. [Licencia](#licencia)

---

## Requisitos
- Go ≥ 1.22
- Ollama instalado y ejecutándose (`ollama serve`).
- Al menos un modelo descargado (`ollama pull llama2`, `ollama pull gemma`, etc.).
- Linux/macOS (Windows WSL debería funcionar).

## Instalación rápida
```bash
# 1. Clonar repo
 git clone https://github.com/tuuser/go-ollama-tui.git && cd go-ollama-tui

# 2. Descargar dependencias
 make deps   # = go mod download

# 3. Ejecutar la TUI
 make run
```

La TUI detectará los modelos disponibles y mostrará un popup de selección al inicio (y con `Ctrl+M`).

## Uso
| Acción | Tecla |
|--------|-------|
| Abrir/cambiar modelo | `Ctrl+M` |
| Enviar mensaje | **Enter** |
| Limpiar chat | `Ctrl+L` |
| Scroll | `↑ ↓` o `PgUp` / `PgDn` |
| Salir | `Ctrl+C` o `Esc` |

La barra inferior muestra el total de tokens recibidos y la velocidad de generación (*tokens/seg*).

### Renderizado Markdown
El contenido se procesa en streaming con [Glamour](https://github.com/charmbracelet/glamour) para soportar listas, tablas, código, etc.

## Arquitectura
```
cmd/tui          → entrypoint principal
internal/
  api/           → cliente HTTP minimal para Ollama (chat & modelos)
  model/         → tipos de dominio (ChatRequest, ChatResponse)
  ui/            → modelo Bubble Tea, estilos Lipgloss, lógica de interacción
  mcp/           → (opcional) utilidades para Model Context Protocol
```
Componentes Bubble Tea utilizados: `textinput`, `viewport`, `spinner`, `list`.

## Desarrollo y mantenimiento
```bash
# Formatear
make fmt      # = go fmt ./...

# Actualizar deps
make update   # = go get -u ./...

# Ejecutar pruebas (pendientes de añadir)
make test     # = go test ./...
```

### Entorno virtual (opcional)
Para quien prefiera aislamiento:
```bash
python -m venv venv && source venv/bin/activate
```

## Hoja de ruta
- [x] Streaming de tokens
- [x] Renderizado Markdown
- [x] Selección dinámica de modelo
- [ ] Soporte MCP estable
- [ ] Exportar historial a archivo
- [ ] Pruebas unitarias + CI

## Licencia
MIT. Ver `LICENSE`.

---

creado por N.A.V.
