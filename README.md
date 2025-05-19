# Stock Insights API

Un servicio para obtener, almacenar y recomendar acciones (stocks) basado en datos de casas de bolsa.

## Tabla de Contenidos

- [Requisitos Previos](#requisitos-previos)
- [Configuración Local](#configuración-local)
- [Endpoints API](#endpoints-api)
- [Infraestructura con Terraform](#infraestructura-con-terraform)
- [Despliegue en GCP](#despliegue-en-gcp)
- [Despliegue con GitHub Actions](#despliegue-con-github-actions)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Variables de Entorno](#variables-de-entorno)
- [Soporte Docker](#soporte-docker)

## Requisitos Previos

- Go 1.18 o superior
- Docker (para ejecutar CockroachDB)
- Git
- Cuenta en Google Cloud Platform (para despliegue en la nube)
- Terraform 1.0+ (para aprovisionamiento de infraestructura)

## Configuración Local

### 1. Clonar el repositorio

```bash
git clone https://github.com/RobertCastro/stock-insights-api
cd stock-insights-api
```

### 2. Instalar dependencias

```bash
go mod tidy
```

### 3. Iniciar CockroachDB con Docker

```bash
docker run -d --name=roach --hostname=roach -p 26257:26257 -p 8080:8080 cockroachdb/cockroach:latest start-single-node --insecure
```

### 4. Crear la base de datos

```bash
docker exec -it roach ./cockroach sql --insecure -e "CREATE DATABASE IF NOT EXISTS stockdb;"
```

### 5. Configurar variables de entorno

Crea un archivo `.env` en la raíz del proyecto y configura las variables necesarias:

```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=26257
DB_USER=root
DB_PASSWORD=
DB_NAME=stockdb
DB_SSL_MODE=disable
```

### 6. Ejecutar la aplicación

```bash
go run cmd/api/main.go
```

## Endpoints API

El servicio expone los siguientes endpoints:

- `GET /api/v1/stocks` - Lista todas las acciones con opciones de filtrado por ticker, brokerage, rating y ordenamiento
- `GET /api/v1/stocks/{ticker}` - Obtiene detalles de una acción específica
- `GET /api/v1/recommendations` - Obtiene recomendaciones de acciones
- `POST /api/v1/sync` - Sincroniza datos desde la API externa
- `GET /health` - Verifica el estado del servicio

### API Externa

Ejemplo de respuesta:
```json
{
    "items": [
        {
            "ticker": "BSBR",
            "target_from": "$4.20",
            "target_to": "$4.70",
            "company": "Banco Santander (Brasil)",
            "action": "upgraded by",
            "brokerage": "The Goldman Sachs Group",
            "rating_from": "Sell",
            "rating_to": "Neutral",
            "time": "2025-01-13T00:30:05.813548892Z"
        }
    ],
    "next_page": "AZEK"
}
```

## Infraestructura con Terraform

Este proyecto incluye configuración Terraform para desplegar la infraestructura necesaria en Google Cloud Platform.

### Componentes de infraestructura

- VPC y subred privada
- Cluster GKE (Google Kubernetes Engine)
- Router y NAT para conexiones salientes
- Cuenta de servicio con permisos adecuados

### Configuración de Terraform

1. Navega al directorio de Terraform:

```bash
cd terraform
```

2. Crea un archivo `terraform.tfvars` basado en el ejemplo:

```bash
cp terraform.tfvars.example terraform.tfvars
```

3. Edita el archivo `terraform.tfvars` para configurar tus variables específicas:

```
project_id = "tu-proyecto-gcp"
region     = "us-central1"
zone       = "us-central1-a"

backend_repo  = "https://github.com/RobertCastro/stock-recommendation-service"
backend_image = "golang-backend"

gke_node_count   = 2
gke_machine_type = "e2-standard-2"
```

4. Inicializa Terraform:

```bash
terraform init
```

5. Verifica el plan de Terraform:

```bash
terraform plan
```

6. Aplica la configuración:

```bash
terraform apply
```

## Despliegue en GCP

### Requisitos para el despliegue

- Cuenta de GCP con facturación habilitada
- Proyecto de GCP creado
- APIs habilitadas:
  - Compute Engine API
  - Kubernetes Engine API
  - Artifact Registry API
  - Cloud Build API

## Despliegue con GitHub Actions

El proyecto incluye un flujo de trabajo de GitHub Actions que automatiza el proceso de construcción y despliegue a Google Cloud Platform. El flujo se activa cuando se realiza un push a la rama `main` o a ramas con el prefijo `feat/`.

#### Requisitos previos para GitHub Actions

1. Configura los siguientes secrets en tu repositorio GitHub:
   - `GCP_SA_KEY`: Archivo JSON de credenciales de una cuenta de servicio con permisos para GKE y Artifact Registry
   - `GCP_PROJECT_ID`: ID de tu proyecto de Google Cloud
   - `GCP_ZONE`: Zona donde está desplegado tu cluster de GKE

#### Proceso de despliegue automático

El flujo de trabajo (`deploy-backend.yml`) realiza las siguientes acciones:

1. Autentica con Google Cloud usando la cuenta de servicio
2. Configura el SDK de Google Cloud
3. Instala el plugin de autenticación GKE
4. Configura Docker para usar Artifact Registry
5. Construye y sube la imagen Docker a Artifact Registry
6. Obtiene el digest de la imagen para un despliegue inmutable
7. Despliega CockroachDB si no está ya desplegado
8. Despliega o actualiza la aplicación backend con la nueva imagen

Para verificar el estado del despliegue, puedes revisar:
- La pestaña "Actions" en tu repositorio de GitHub
- El estado de los recursos en GKE:
  ```bash
  kubectl get pods
  kubectl get services
  ```

#### Capturas de pantalla

![image](https://github.com/user-attachments/assets/e8eb4bd3-89ae-4561-826c-ee778afc3a85)
![image](https://github.com/user-attachments/assets/162f310e-8b4a-47c2-b7f5-5f22e327158c)




### Despliegue manual (alternativa)

Para realizar el despliegue manualmente, sigue estos pasos:

1. Autentícate con Google Cloud:

```bash
gcloud auth login
gcloud config set project [TU_PROYECTO_ID]
```

2. Construye la imagen de Docker:

```bash
docker build -t us-central1-docker.pkg.dev/[TU_PROYECTO_ID]/docker-images/golang-backend:latest .
```

3. Configura Docker para usar Artifact Registry:

```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

4. Sube la imagen a Artifact Registry:

```bash
docker push us-central1-docker.pkg.dev/[TU_PROYECTO_ID]/docker-images/golang-backend:latest
```

5. Conecta con el cluster de Kubernetes:

```bash
gcloud container clusters get-credentials backend-cluster --zone [TU_ZONA] --project [TU_PROYECTO_ID]
```

6. Aplica las configuraciones de Kubernetes:

```bash
kubectl apply -f kubernetes/cockroachdb-k8s.yaml
kubectl apply -f kubernetes/deployment.yaml
kubectl apply -f kubernetes/ingress.yaml
```

7. Verifica el despliegue:

```bash
kubectl get pods
kubectl get services
```

## Variables de Entorno

| Variable      | Descripción                        | Valor por defecto |
|---------------|------------------------------------|-------------------|
| SERVER_PORT   | Puerto del servidor                | 8080              |
| DB_HOST       | Host de la base de datos           | localhost         |
| DB_PORT       | Puerto de la base de datos         | 26257             |
| DB_USER       | Usuario de la base de datos        | root              |
| DB_PASSWORD   | Contraseña de la base de datos     |                   |
| DB_NAME       | Nombre de la base de datos         | stockdb           |
| DB_SSL_MODE   | Modo SSL para la conexión          | disable           |
| API_KEY       | Token de autenticación para la API externa |           |

## Soporte Docker

El proyecto incluye un Dockerfile para facilitar la creación de contenedores. Para construir y ejecutar la aplicación usando Docker:

### Construir la imagen

```bash
docker build -t stock-service:latest .
```

### Ejecutar el contenedor

```bash
docker run -p 8080:8080 --env-file .env stock-service:latest
```

### Docker Compose (opcional)

Si prefieres usar Docker Compose para manejar tanto la aplicación como CockroachDB, puedes crear un archivo `docker-compose.yml`:

```yaml
version: '3'

services:
  roach:
    image: cockroachdb/cockroach:latest
    hostname: roach
    command: start-single-node --insecure
    ports:
      - "26257:26257"
      - "8081:8080"
    volumes:
      - cockroach-data:/cockroach/cockroach-data

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_PORT=8080
      - DB_HOST=roach
      - DB_PORT=26257
      - DB_USER=root
      - DB_PASSWORD=
      - DB_NAME=stockdb
      - DB_SSL_MODE=disable
    depends_on:
      - roach

volumes:
  cockroach-data:
```

Y ejecutarlo con:

```bash
docker-compose up -d
```

## Notas Importantes

- La aplicación utiliza una arquitectura hexagonal para separar las responsabilidades y facilitar las pruebas.
- La sincronización con la API externa se puede configurar para ejecutarse automáticamente con un cron job.
- Para entornos de producción, habilita SSL para la conexión a la base de datos.
- El token de autenticación para la API externa debe mantenerse seguro, preferiblemente usando un gestor de secretos en entornos de producción.
