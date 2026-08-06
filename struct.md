# 🚀 Go Backend Project Template & Architecture Guide (`struct.md`)

คู่มือและโครงสร้างมาตรฐานสำหรับการเริ่มต้นโปรเจกต์ใหม่ (Boilerplate / Starter Template) ด้วย **Go (Golang)** โดยใช้ **Fiber Framework** ตามแนวคิด **Hexagonal Architecture (Ports & Adapters)** และ **Modular Monolith**

---

## 📌 1. ภาพรวมสถาปัตยกรรม (Architecture Concepts)

โปรเจกต์จะถูกจัดโครงสร้างให้แต่ละ **Module** แยกขอบเขตการทำงานกันอย่างชัดเจน โดยใช้แนวคิด **Hexagonal Architecture**:

* **Ports (Interfaces)**: กำหนดข้อตกลง (Contracts) ของ Service และ Repository เพื่อให้เกิด Loose Coupling และง่ายต่อการเขียน Unit Test
* **Adapters (Implementations)**:
  * **Inbound Adapter**: HTTP Handlers (Fiber Handlers)
  * **Outbound Adapter**: Repositories (GORM/SQL, MongoDB, External API Client)
* **Core Logic**: Business Logic อยู่ใน `services/` ไม่ขึ้นกับ Web Framework หรือ Database

---

## 🗂️ 2. โครงสร้างโฟลเดอร์มาตรฐาน (Standard Directory Tree)

```text
<project-name>/
├── app/                  # Application Lifecycle, Dependency Injection & Server Setup
├── cmd/                  # Entry points สำหรับการรันแอปพลิเคชัน
│   └── server/           # Main entrypoint (main.go) สำหรับรัน HTTP Server
├── config/               # Struct และ Loader สำหรับอ่านไฟล์ Configuration (YAML/Env)
├── model/                # Database Entities & Schemas (GORM Models / DB Tables)
├── module/               # Business Feature Modules (Hexagonal Structure)
│   ├── [feature_a]/      # โมดูลระบบ เช่น auth, user, product, order
│   │   ├── dto/          # Data Transfer Objects (Request/Response validation)
│   │   ├── handler/      # Fiber HTTP Handlers (Controllers)
│   │   ├── mapper/       # Data Mappers (DTO ↔ Entity ↔ Model)
│   │   ├── middleware/   # Module-specific Middlewares (Auth, Roles check)
│   │   ├── ports/        # Interfaces สำหรับ Service และ Repository
│   │   ├── repositories/ # Data Access Layer (GORM/Database Implementation)
│   │   ├── services/     # Business Logic Layer
│   │   ├── tests/        # Unit & Integration Tests
│   │   └── [feature_a].go# Module Router & Dependency Wiring
│   └── module.go         # Common Module Interface & Registry Manager
├── pkg/                  # Shared Utility Packages & Global Clients
│   ├── bcrypt/           # Password Hashing & Verification
│   ├── constant/         # Global Constants & Error Codes
│   ├── database/         # DB Client Drivers (MariaDB, PostgreSQL, MongoDB, Redis)
│   ├── jwt/              # JWT Token Generation & Verification
│   ├── response/         # Standardized REST API Response Formatter
│   ├── tracing/          # OpenTelemetry / Distributed Tracing
│   └── util/             # Utility Helpers (Validators, Date/String Formatters)
├── docs/                 # Documentation & API Specs (Swagger/OpenAPI)
├── k8s/                  # Kubernetes Deployment Manifests
├── .air.toml             # Configuration สำหรับ Live Reloading (Air tool)
├── .env.example          # ตัวอย่างไฟล์ Environment Variables
├── config.yml            # ไฟล์ Configuration หลักสำหรับระบบ
├── Dockerfile            # Docker Image Build file
├── docker-compose.yml    # Local Development Environment (Databases, Cache, Services)
├── Makefile              # Shell Command Automation (build, run, test, dev)
├── go.mod / go.sum       # Go Modules Dependency Manager
└── README.md             # เอกสารเริ่มต้นและการใช้งานโปรเจกต์
```

---

## 🔍 3. หน้าที่และความรับผิดชอบของแต่ละส่วน

### 3.1 🚀 `app/` (Application Lifecycle & DI)
* **`app.go`**: ทำหน้าที่เป็นศูนย์กลาง **Dependency Injection (DI)** เชื่อมต่อ Database Clients (GORM, Mongo, Redis) และทำ Wiring ลงทะเบียนโมดูลต่างๆ เข้าด้วยกัน
* **`server.go`**: ตั้งค่า Fiber Engine, ลงทะเบียน Global Middleware (CORS, Recover, Logger, Rate Limit) และจัดการ **Graceful Shutdown**

### 3.2 ⚡ `cmd/` (Application Entrypoints)
* **`cmd/server/main.go`**: จุดเริ่มต้นการรันโปรแกรม (`go run cmd/server/main.go start`) อ่านค่า Config แล้วส่งต่อให้ `app/` ทำงาน

### 3.3 ⚙️ `config/` (Configuration Management)
* **`config.go`**: อ่านค่าจากไฟล์ `config.yml` หรือ Environment Variables แล้ว Map เข้ากับ Go Structs เพื่อให้เรียกใช้ค่าคอนฟิกต่าง ๆ ได้อย่างมี Type-Safety (เช่น Database URL, JWT Secrets, Server Port)

### 3.4 🗄️ `model/` (Domain Entities & Data Schemas)
* นิยาม Struct ของตารางในฐานข้อมูล (GORM Entities) เช่น `User`, `Role`, `Product`
* เก็บโครงสร้างที่เป็น Core Data Structures ของระบบ

### 3.5 🧱 `module/` (Feature Modules Layout)
ทุกโมดูลภายใต้ `module/` จะแบ่งชั้นการทำงานออกเป็น 6 Layer หลัก:

1. **`ports/`**: นิยาม Interface ข้อตกลง เช่น `UserService` และ `UserRepository`
2. **`dto/`**: Struct สำหรับรับ Request (พร้อม Validation tags) และส่ง Response
3. **`services/`**: เขียน Business Logic หลัก ห้ามเรียกใช้ SQL/Database โดยตรง แต่ให้เรียกผ่าน Repository Interface
4. **`repositories/`**: เขียนโค้ดติดต่อฐานข้อมูลจริง (GORM, SQL Queries, Mongo Client)
5. **`handler/`**: Fiber HTTP Controller รับ Request ป้อนเข้า DTO แล้วเรียก Service จากนั้นคืนค่าด้วย Standard Response
6. **`middleware/`**: Middleware เฉพาะของโมดูล เช่น ตรวจสอบสิทธิ์การเข้าถึงข้อมูลเฉพาะส่วน

### 3.6 🛠️ `pkg/` (Shared Packages)
แพ็กเกจกลางที่เป็นอิสระ สามารถนำไป Re-use ในโมดูลอื่น ๆ หรือโปรเจกต์อื่นได้:
* **`pkg/response/`**: Format การตอบกลับแบบมาตรฐาน เช่น `{ "code": 200, "message": "success", "data": ... }`
* **`pkg/jwt/`**: การสร้างและยืนยันตัวตนผ่าน JWT Access/Refresh Token
* **`pkg/bcrypt/`**: การเข้ารหัสผ่าน
* **`pkg/constant/`**: รวม System Error Codes และ Status Text กลาง
* **`pkg/util/`**: เครื่องมือช่วยประมวลผลคำ, วันที่, หรือ Validation ต่าง ๆ

---

## 🛠️ 4. ขั้นตอนการเพิ่ม Feature Module ใหม่ (Step-by-Step)

เมื่อต้องการเพิ่มฟีเจอร์ใหม่ ให้ดำเนินการตามลำดับดังนี้:

```text
Step 1: สร้าง Model       -->  model/product.go
Step 2: กำหนด Ports       -->  module/product/ports/product_ports.go
Step 3: ทำ Repository     -->  module/product/repositories/product_repository.go
Step 4: เขียน Business Logic --> module/product/services/product_service.go
Step 5: ทำ HTTP Handler   -->  module/product/handler/product_handler.go
Step 6: สร้าง DTO          -->  module/product/dto/product_dto.go
Step 7: Register Module   -->  module/product/product.go & app/app.go
```

---

## 📐 5. มาตรฐาน JSON API Response Format (`pkg/response`)

### Response สำเร็จ (Success):
```json
{
  "code": 200,
  "status": "success",
  "message": "Operation completed successfully",
  "data": {}
}
```

### Response ผิดพลาด (Error):
```json
{
  "code": 400,
  "status": "error",
  "message": "Invalid request payload",
  "errors": [
    {
      "field": "email",
      "message": "email is a required field"
    }
  ]
}
```

---

## ⚙️ 6. ตัวอย่างคำสั่งพื้นฐานใน `Makefile`

```makefile
# Run with Hot-reload using Air
run:
	air

# Run server directly
dev:
	go run cmd/server/main.go start

# Run all unit tests
test:
	go test -v -coverprofile=coverage.out ./...

# Build binary for Production
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/server ./cmd/server/main.go
```
