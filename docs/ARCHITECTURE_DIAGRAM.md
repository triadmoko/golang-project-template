# 🏗️ Architecture Diagram

## Feature-Based Clean Architecture

```mermaid
graph TB
    subgraph "External Layer"
        UI[Web UI / Mobile App]
        API[External APIs]
    end

    subgraph "Delivery Layer"
        subgraph "Auth Feature"
            AH[Auth Handler]
            AD[Auth DTOs]
        end

        subgraph "User Feature"
            UH[User Handler]
            UD[User DTOs]
        end

        subgraph "Product Feature"
            PH[Product Handler]
            PD[Product DTOs]
        end

        MW[Middleware]
        R[Router]
    end

    subgraph "Use Case Layer"
        subgraph "Auth Feature"
            AU[Auth Use Case]
        end

        subgraph "User Feature"
            UU[User Use Case]
        end

        subgraph "Product Feature"
            PU[Product Use Case]
        end
    end

    subgraph "Domain Layer"
        subgraph "Auth Feature"
            AE[User Entity]
            AR[User Repository Interface]
            AS[Auth Service Interface]
        end

        subgraph "User Feature"
            UE[User Entity]
            UR[User Repository Interface]
        end

        subgraph "Product Feature"
            PE[Product Entity]
            PR[Product Repository Interface]
        end

        SE[Shared Errors]
    end

    subgraph "Infrastructure Layer"
        subgraph "Auth Feature"
            AIR[User Repository Impl]
            AIS[Auth Service Impl]
        end

        subgraph "User Feature"
            UIR[User Repository Impl]
        end

        subgraph "Product Feature"
            PIR[Product Repository Impl]
        end

        DB[(PostgreSQL Database)]
        EXT[External Services]
    end

    UI --> R
    API --> R
    R --> MW
    MW --> AH
    MW --> UH
    MW --> PH

    AH --> AU
    UH --> UU
    PH --> PU

    AU --> AE
    AU --> AR
    AU --> AS
    UU --> UE
    UU --> UR
    PU --> PE
    PU --> PR

    AIR --> AR
    AIS --> AS
    UIR --> UR
    PIR --> PR

    AIR --> DB
    UIR --> DB
    PIR --> DB
    AIS --> EXT
```

## Layer Dependencies

```mermaid
graph LR
    subgraph "Clean Architecture Layers"
        D[Delivery Layer]
        U[Use Case Layer]
        D2[Domain Layer]
        I[Infrastructure Layer]
    end

    D --> U
    U --> D2
    I --> D2

    style D fill:#e1f5fe
    style U fill:#f3e5f5
    style D2 fill:#e8f5e8
    style I fill:#fff3e0
```

## Feature Isolation

```mermaid
graph TB
    subgraph "Auth Feature"
        AD[Auth Domain]
        AU[Auth Use Case]
        AI[Auth Infrastructure]
        AH[Auth Handler]
    end

    subgraph "User Feature"
        UD[User Domain]
        UU[User Use Case]
        UI[User Infrastructure]
        UH[User Handler]
    end

    subgraph "Product Feature"
        PD[Product Domain]
        PU[Product Use Case]
        PI[Product Infrastructure]
        PH[Product Handler]
    end

    subgraph "Shared Components"
        SC[Shared Domain]
        SM[Shared Middleware]
        SR[Shared Router]
    end

    AD -.-> SC
    UD -.-> SC
    PD -.-> SC

    AH --> SM
    UH --> SM
    PH --> SM

    AH --> SR
    UH --> SR
    PH --> SR
```

## Data Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router
    participant H as Handler
    participant U as Use Case
    participant Rep as Repository
    participant DB as Database

    C->>R: HTTP Request
    R->>H: Route to Handler
    H->>U: Call Use Case
    U->>Rep: Call Repository
    Rep->>DB: Query Database
    DB-->>Rep: Return Data
    Rep-->>U: Return Entity
    U-->>H: Return Result
    H-->>R: Return Response
    R-->>C: HTTP Response
```

## Feature Dependencies

```mermaid
graph TD
    subgraph "Core Features"
        A[Auth Feature]
        U[User Feature]
        P[Product Feature]
    end

    subgraph "Future Features"
        O[Order Feature]
        PAY[Payment Feature]
        INV[Inventory Feature]
    end

    subgraph "Shared"
        S[Shared Components]
    end

    A --> S
    U --> S
    P --> S

    O -.-> A
    O -.-> P
    PAY -.-> A
    PAY -.-> O
    INV -.-> P

    style A fill:#ffebee
    style U fill:#e8f5e8
    style P fill:#e3f2fd
    style S fill:#fff3e0
```

## Microservices Migration Path

```mermaid
graph TB
    subgraph "Current: Modular Monolith"
        MM[Modular Monolith]
        AF[Auth Feature]
        UF[User Feature]
        PF[Product Feature]
    end

    subgraph "Future: Microservices"
        AS[Auth Service]
        US[User Service]
        PS[Product Service]
        GS[API Gateway]
    end

    MM --> AS
    MM --> US
    MM --> PS

    AF --> AS
    UF --> US
    PF --> PS

    GS --> AS
    GS --> US
    GS --> PS

    style MM fill:#e1f5fe
    style AS fill:#ffebee
    style US fill:#e8f5e8
    style PS fill:#e3f2fd
    style GS fill:#fff3e0
```

## Key Benefits

1. **🔒 Isolation**: Each feature is self-contained
2. **🔄 Scalability**: Easy to scale individual features
3. **🧪 Testability**: Independent testing per feature
4. **👥 Team Work**: Teams can work on different features
5. **🚀 Migration**: Easy path to microservices
6. **🛠️ Maintenance**: Easier to maintain and debug
