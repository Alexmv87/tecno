# API de tickets

API REST escrita en Go para crear, consultar, actualizar y eliminar tickets. Usa PostgreSQL para persistir los datos y puede ejecutarse localmente o con Docker Compose.

## Requisitos

- Go 1.27.4 o compatible con la versión indicada en `go.mod`.
- PostgreSQL, si ejecutas la API directamente en tu equipo.
- Docker/Podman y Docker Compose, para levantar la API y la base de datos en contenedores.

## Configuración

La aplicación carga las variables del sistema y, si existe, el archivo `.env` en la raíz del proyecto. Genere un `.env` en la ruta del proyecto con los siguientes valores cambiando el DB_USER y DB_PASSWORD por los valores reales:

```env
PORT=8080
DB_USER=usuarioparabasededatos
DB_PASSWORD=indicarpasswordseguro
DB_HOST=api-db
DB_NAME=tecno_db
DB_PORT=5432
```

## Ejecución local

1. Inicia PostgreSQL y crea la base de datos configurada en `DATABASE_URI`.
2. Configura el archivo `.env` con el host `localhost`.
3. Desde la raíz del proyecto, ejecuta:

   ```sh
   go run .
   ```

Al iniciar, la aplicación crea automáticamente la tabla `tickets` si aún no existe. La API escucha en el puerto indicado por `PORT`.

## Ejecución con Docker Compose

Configura el `.env` como se indica arriba y ejecuta:

```sh
docker compose up --build
```

La API queda disponible en `http://localhost:8080` y PostgreSQL se publica localmente en el puerto `5432`. Para detener los servicios:

```sh
docker compose down
```

## Endpoints

La base de las rutas de tickets es `/tickets`.

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/health` | Comprueba la conexión con PostgreSQL. Devuelve `200` si está disponible y `503` si no lo está. |
| `GET` | `/tickets/` | Lista los tickets. |
| `GET` | `/tickets/{id}` | Obtiene un ticket por su identificador. |
| `POST` | `/tickets/` | Crea un ticket. |
| `PUT` | `/tickets/{id}` | Actualiza el estado de un ticket. |
| `DELETE` | `/tickets/{id}` | Elimina un ticket. |

Ejemplo de cuerpo JSON para crear un ticket:

```json
{
  "cliente": "Ana",
  "asunto": "Problema con el pedido"
}
```

Para actualizar, el cuerpo JSON utiliza el campo `status`:

```json
{
  "status": "resuelto"
}
```

Los tickets se devuelven como JSON con los campos `id`, `cliente`, `asunto`, `status` y `creado_en`. El estado inicial en la base de datos es `abierto`.

## Persistencia

La tabla `tickets` se crea automáticamente con estas columnas:

| Columna | Tipo | Descripción |
|---|---|---|
| `id` | `SERIAL` | Identificador autoincremental y clave primaria. |
| `cliente` | `VARCHAR(100)` | Cliente que reporta el ticket. |
| `asunto` | `VARCHAR(200)` | Asunto del ticket. |
| `estado` | `VARCHAR(20)` | Estado inicial `abierto`. |
| `creado_en` | `TIMESTAMP` | Fecha y hora de creación. |
