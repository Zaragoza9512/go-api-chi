# 🏗️ Arquitectura del Sistema

## Diagrama General
```mermaid
graph TB
    subgraph "Cliente"
        Client[👤 Usuario/Cliente]
    end
    
    subgraph "AWS Cloud"
        subgraph "Compute"
            EC2[🖥️ EC2 Instance<br/>t3.micro<br/>Amazon Linux 2023]
        end
        
        subgraph "Database"
            RDS[(🗄️ RDS PostgreSQL 17.2<br/>db.t4g.micro<br/>20GB SSD)]
        end
        
        subgraph "Terraform State"
            S3[📦 S3 Bucket<br/>terraform-state<br/>Versioning enabled]
            DynamoDB[🔒 DynamoDB Table<br/>terraform-locks<br/>State Locking]
        end
        
        subgraph "Security"
            SG1[🛡️ Security Group<br/>EC2<br/>Ports: 22, 8080]
            SG2[🛡️ Security Group<br/>RDS<br/>Port: 5432]
        end
    end
    
    subgraph "CI/CD"
        GitHub[📦 GitHub<br/>Repository]
        Actions[⚙️ GitHub Actions<br/>CI/CD Pipeline]
        DockerHub[🐳 Docker Hub<br/>Image Registry]
    end
    
    subgraph "Monitoring"
        Prometheus[📊 Prometheus<br/>Metrics Collection]
        Grafana[📈 Grafana<br/>Dashboards]
    end
    
    Client -->|HTTPS| EC2
    EC2 -->|SQL| RDS
    EC2 -.->|Metrics| Prometheus
    Prometheus -.->|Query| Grafana
    
    GitHub -->|Push| Actions
    Actions -->|Build & Test| Actions
    Actions -->|Push Image| DockerHub
    DockerHub -.->|Pull Image| EC2
    
    EC2 ---|Protected by| SG1
    RDS ---|Protected by| SG2
    
    S3 -.->|Stores| TerraformState[terraform.tfstate]
    DynamoDB -.->|Locks| TerraformState
    
    style EC2 fill:#ff9900
    style RDS fill:#336791
    style S3 fill:#569A31
    style DynamoDB fill:#4053D6
    style Prometheus fill:#E6522C
    style Grafana fill:#F46800
```

## Flujo de Despliegue
```mermaid
sequenceDiagram
    participant Dev as 👨‍💻 Developer
    participant Git as GitHub
    participant GHA as GitHub Actions
    participant DH as Docker Hub
    participant EC2 as AWS EC2
    participant RDS as AWS RDS
    
    Dev->>Git: 1. git push
    Git->>GHA: 2. Trigger workflow
    GHA->>GHA: 3. Run tests
    GHA->>GHA: 4. Build Docker image
    GHA->>DH: 5. Push image
    Note over GHA,DH: Image: zaragoza95/go-api-chi:latest
    DH-->>EC2: 6. Pull image
    EC2->>EC2: 7. Deploy container
    EC2->>RDS: 8. Connect to DB
    EC2-->>Dev: 9. API Ready ✅
```

## Infraestructura (Terraform)
```mermaid
graph LR
    subgraph "Developer Machine"
        TF[💻 Terraform CLI]
    end
    
    subgraph "Remote State"
        S3[📦 S3 Bucket<br/>State Storage]
        DB[(🔒 DynamoDB<br/>State Lock)]
    end
    
    subgraph "AWS Resources"
        EC2[🖥️ EC2]
        RDS[🗄️ RDS]
        SG[🛡️ Security Groups]
    end
    
    TF -->|1. terraform init| S3
    TF -->|2. Acquire lock| DB
    TF -->|3. terraform apply| EC2
    TF -->|4. terraform apply| RDS
    TF -->|5. terraform apply| SG
    TF -->|6. Update state| S3
    TF -->|7. Release lock| DB
    
    style S3 fill:#569A31
    style DB fill:#4053D6
    style EC2 fill:#ff9900
    style RDS fill:#336791
```

## Stack de Monitoreo
```mermaid
graph TB
    subgraph "Go API"
        API[🚀 Chi Router<br/>Port 8080]
        Metrics[📊 /metrics Endpoint<br/>Prometheus format]
    end
    
    subgraph "Monitoring Stack"
        Prom[📈 Prometheus<br/>Port 9090]
        Graf[📊 Grafana<br/>Port 3000]
    end
    
    API -->|Exposes| Metrics
    Prom -->|Scrapes every 30s| Metrics
    Graf -->|Queries| Prom
    
    User[👤 DevOps Team] -->|Views| Graf
    
    style API fill:#00ADD8
    style Prom fill:#E6522C
    style Graf fill:#F46800
```

## Flujo de Request
```mermaid
sequenceDiagram
    participant C as 👤 Client
    participant LB as Load Balancer
    participant API as Go API
    participant JWT as JWT Middleware
    participant DB as PostgreSQL
    participant M as Prometheus
    
    C->>LB: 1. HTTPS Request
    LB->>API: 2. Forward to API
    API->>JWT: 3. Validate Token
    
    alt Token Valid
        JWT->>API: 4. Authorized
        API->>DB: 5. Query Data
        DB-->>API: 6. Return Data
        API-->>C: 7. JSON Response (200)
        API->>M: 8. Record Metrics
    else Token Invalid
        JWT-->>C: 7. Unauthorized (401)
        API->>M: 8. Record Error
    end
```

## Componentes Principales

### Backend (Go)
- **Framework:** Chi Router v5
- **Database:** PostgreSQL driver (lib/pq)
- **Auth:** JWT (golang-jwt)
- **Metrics:** Prometheus client

### Infrastructure (AWS)
- **Compute:** EC2 t3.micro
- **Database:** RDS PostgreSQL 17.2
- **Storage:** S3 (Terraform state)
- **Locking:** DynamoDB (Terraform)

### DevOps
- **IaC:** Terraform 1.14+
- **CI/CD:** GitHub Actions
- **Containers:** Docker + Docker Compose
- **Registry:** Docker Hub

### Monitoring
- **Metrics:** Prometheus
- **Visualization:** Grafana
- **Method:** RED (Rate, Errors, Duration)

---

## Seguridad

### Network
- ✅ RDS en subnet privada (no public IP)
- ✅ Security Groups con least privilege
- ✅ HTTPS en producción

### Application
- ✅ JWT para autenticación
- ✅ Middleware de autorización
- ✅ Validación de inputs

### Infrastructure
- ✅ Terraform state encriptado en S3
- ✅ Secrets en variables de entorno
- ✅ .gitignore para credenciales

---

## Escalabilidad

### Horizontal
- Auto Scaling Groups (futuro)
- Load Balancer (futuro)
- Multiple AZs (futuro)

### Vertical
- Ajustar instance types en Terraform
- Variables parametrizadas
- Sin downtime con Blue/Green

### Database
- RDS Read Replicas (futuro)
- Connection pooling
- Query optimization