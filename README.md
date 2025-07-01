# GOX 🚀

[![CI](https://github.com/whitejokeer/gox/workflows/CI/badge.svg)](https://github.com/whitejokeer/gox/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/whitejokeer/gox)](https://goreportcard.com/report/github.com/whitejokeer/gox)
[![codecov](https://codecov.io/gh/whitejokeer/gox/branch/main/graph/badge.svg)](https://codecov.io/gh/whitejokeer/gox)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.24-blue)](https://golang.org/)
[![Release](https://img.shields.io/github/release/whitejokeer/gox.svg)](https://github.com/whitejokeer/gox/releases)

GOX es un framework web moderno que unifica Go, HTMX y CSS en componentes de archivo único (.gox), ofreciendo una experiencia de desarrollo similar a Vue/Svelte pero con la simplicidad y rendimiento del server-side rendering.

## ✨ Características

- 🎯 **Single File Components**: Todo en un archivo `.gox` - template, lógica y estilos
- ⚡ **Server-Side Rendering**: Rendimiento nativo de Go con HTML generado en servidor
- 🔄 **HTMX Integration**: Interactividad moderna sin JavaScript complejo
- 🎨 **Scoped CSS**: Estilos automáticamente encapsulados por componente
- 🔥 **Hot Reload**: Desarrollo rápido con recarga automática
- 📦 **Zero Dependencies**: Solo Go puro, sin runtime JavaScript

## 🚀 Instalación

```bash
# Instalar desde código fuente
go install github.com/whitejokeer/gox/cmd/gox@latest

# O clonar y compilar
git clone https://github.com/whitejokeer/gox.git
cd gox
make install
```

## 📖 Uso Rápido

### 1. Crear un componente

```gox
<!-- hello.gox -->
<template>
  <div class="hello-component">
    <h1>{{ .Title }}</h1>
    <button hx-post="/click" hx-target="#result">
      Click me!
    </button>
    <div id="result"></div>
  </div>
</template>

<script>
package main

import "net/http"

type HelloProps struct {
  Title string `json:"title"`
}

func (h HelloProps) HandleClick(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("<p>¡Hola desde GOX!</p>"))
}
</script>

<style>
.hello-component {
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.hello-component h1 {
  color: #333;
  margin-bottom: 1rem;
}
</style>
```

### 2. Compilar componentes

```bash
# Compilar un componente específico
gox build hello.gox

# Compilar todos los componentes en un directorio
gox build src/

# Modo desarrollo con hot reload
gox dev

# Vigilar cambios
gox watch src/
```

## 📁 Estructura del Proyecto

```
gox/
├── cmd/gox/           # CLI principal
├── internal/          # Paquetes internos
│   ├── parser/        # Analizador de archivos .gox
│   ├── compiler/      # Compilador de componentes
│   └── watcher/       # Vigilancia de archivos
├── pkg/               # API pública
├── templates/         # Plantillas de componentes
├── examples/          # Ejemplos de uso
└── .github/workflows/ # CI/CD
```

## 🛠️ Desarrollo

```bash
# Configurar entorno de desarrollo
make dev-setup

# Ejecutar tests
make test

# Linting
make lint

# Compilar
make build

# Ver todos los comandos disponibles
make help
```

## 🧪 Testing

```bash
# Ejecutar todos los tests
make test

# Test con coverage
make test-coverage

# Ejecutar checks de calidad
make check
```

## 📚 Documentación

- [Guía de Inicio](docs/getting-started.md)
- [Sintaxis de Componentes](docs/component-syntax.md)
- [Integración HTMX](docs/htmx-integration.md)
- [API Reference](docs/api-reference.md)

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'Add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia Apache 2.0. Ver [LICENSE](LICENSE) para más detalles.

## 🙋‍♂️ Soporte

- 📖 [Documentación](docs/)
- 🐛 [Issues](https://github.com/whitejokeer/gox/issues)
- 💬 [Discussions](https://github.com/whitejokeer/gox/discussions)

---

Hecho con ❤️ para la comunidad Go
