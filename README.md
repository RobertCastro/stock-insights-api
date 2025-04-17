# Stock Recommendation Service

Un servicio para obtener, almacenar y recomendar acciones (stocks) basado en datos de casas de bolsa.

## Requisitos previos

- Go 1.18 o superior
- Docker (para ejecutar CockroachDB)
- Acceso a internet (para consumir la API externa)

## Instalación

1. Clonar el repositorio:

```bash
git clone https://github.com/RobertCastro/stock-recommendation-service.git
cd stock-recommendation-service
```

2. Instalar dependencias:

```bash
go mod tidy
```

3. Iniciar CockroachDB con Docker:

```bash
docker run -d --name=roach --hostname=roach -p 26257:26257 -p 8080:8080 cockroachdb/cockroach:latest start-single-node --insecure
```

4. Crear la base de datos:

```bash
docker exec -it roach ./cockroach sql --insecure -e "CREATE DATABASE IF NOT EXISTS stockdb;"
```

## Configuración

La aplicación se configura mediante variables de entorno:

| Variable      | Descripción                        | Valor por defecto |
|---------------|------------------------------------|-------------------|
| SERVER_PORT   | Puerto del servidor                | 8080              |
| DB_HOST       | Host de la base de datos           | localhost         |
| DB_PORT       | Puerto de la base de datos         | 26257             |
| DB_USER       | Usuario de la base de datos        | root              |
| DB_PASSWORD   | Contraseña de la base de datos     |                   |
| DB_NAME       | Nombre de la base de datos         | stockdb           |
| DB_SSL_MODE   | Modo SSL para la conexión          | disable           |

## Uso

### Recolección de datos

Para obtener datos de la API y almacenarlos en la base de datos:

```bash
go run cmd/api/main.go
```

## Estructura del Proyecto

```
stock-recommendation-service/
├── cmd/
│   └── api/               # Punto de entrada de la aplicación
├── internal/
│   ├── adapters/
│   │   └── secondary/     # Adaptadores salientes (API, DB)
│   ├── domain/            # Modelos y lógica de dominio
│   └── infrastructure/    # Configuración y utilidades
```

## Endpoints API

La API externa de stocks tiene la siguiente estructura:

```
GET https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list
```

Parámetros de consulta:
- `next_page`: Para paginación

Autenticación:
- Bearer token en el encabezado `Authorization`
