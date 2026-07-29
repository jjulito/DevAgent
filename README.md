# DevAgent CLI

DevAgent es una herramienta de línea de comandos (CLI) desarrollada en Go para asistir en el proceso de desarrollo de software mediante Inteligencia Artificial. Permite realizar consultas sobre el código, ejecutar revisiones automáticas de commits y diffs de Git, indexar y realizar búsquedas semánticas en el proyecto utilizando RAG (Retrieval-Augmented Generation), y exponer funcionalidades mediante el protocolo MCP (Model Context Protocol). Funciona bajo el esquema BYOK (Bring Your Own Key) y soporta modelos locales.

## Cómo funciona

DevAgent interactúa directamente con la terminal del sistema y el repositorio Git local:

1. **Gestión de Proveedores:** Carga la configuración desde archivos de entorno (`.env`), variables de sistema o argumentos de línea de comandos. Se comunica con proveedores remotos (OpenRouter, OpenAI, Google Gemini) o locales (Ollama).
2. **Análisis de Código y Git:** Para la revisión automática, lee el estado del árbol de trabajo de Git (`staged` o `unstaged`) y genera sugerencias orientadas al contexto del proyecto.
3. **Búsqueda Semántica (RAG):** Indexa los archivos fuente del proyecto en una base de datos vectorial (Qdrant). Al realizar búsquedas o preguntas complejas, extrae los fragmentos de código más relevantes y los provee como contexto al modelo.
4. **Servidor MCP:** Implementa el protocolo MCP para actuar como servidor de herramientas que pueden ser invocadas por otros agentes o entornos de desarrollo.

## Utilidades y Comandos

El CLI dispone de las siguientes funcionalidades principales:

- `devagent ask "<pregunta>"`: Consulta al modelo de lenguaje configurado con respuesta en streaming.
- `devagent review`: Genera una revisión automática del código modificado en el repositorio.
  - `--staged`: Analiza únicamente los cambios preparados para commit.
- `devagent index [directorio]`: Indexa los archivos del directorio especificado en la base de datos vectorial Qdrant.
- `devagent search "<consulta>"`: Realiza una búsqueda semántica sobre el código indexado previamente.
- `devagent serve`: Inicia el servidor MCP (puerto por defecto 3000) para integración con agentes externos.
- `devagent config`: Muestra la configuración actual y valida la conectividad con el proveedor activo.

### OPCIONES GLOBALES

- `-p, --provider <nombre>`: Define el proveedor de LLM a utilizar (`openrouter`, `ollama`, `openai`, `gemini`).
- `-m, --model <nombre>`: Define el modelo específico a invocar.
- `-v, --verbose`: Activa el modo detallado de depuración.
- `-c, --config <ruta>`: Especifica la ruta a un archivo de configuración.

## Requisitos y Setup

### Requisitos Previos

- Go 1.21 o superior.
- Docker y Docker Compose (necesarios si se utiliza el almacenamiento vectorial Qdrant para RAG).

### Instalación Local

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/jjulito/DevAgent.git
   cd DevAgent
   ```

2. Configurar el archivo de entorno:
   ```bash
   cp .env.example .env
   ```
   Edite el archivo `.env` para asignar su proveedor y las claves de API correspondientes (`OPENROUTER_API_KEY`, `OPENAI_API_KEY`, `GEMINI_API_KEY`, etc.).

3. Compilar el proyecto:
   ```bash
   make build
   ```
   O bien ejecutar el script de instalación automática:
   ```bash
   bash scripts/setup.sh
   ```

4. Verificar la instalación:
   ```bash
   ./devagent config
   ```

### Ejecución con Docker

Para iniciar la infraestructura de soporte (como la base de datos vectorial Qdrant):

```bash
make docker-up
```

Para detener los servicios de Docker:

```bash
make docker-down
```
