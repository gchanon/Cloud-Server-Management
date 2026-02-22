# Cloud-Server-Management

## Prerequisites

- [Go](https://go.dev/dl/) 1.24+
- Mocked Infrastructure Microservice running on `http://localhost:8081`

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/golf/Cloud-Server-Management.git
cd Cloud-Server-Management
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the backend server (port 8080)

```bash
go run services/cloudMgmt/main.go
```

### 4. Run the frontend server (port 3000)

```bash
cd frontend
go run main.go
```

เปิด [http://localhost:3000](http://localhost:3000) ใน browser (เพื่อให้สามารถ test CORS ได้)

## API Endpoints

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/auth/login` | Login and receive JWT cookie | No |
| GET | `/servers` | Get all provisioned servers | Yes |
| POST | `/servers` | Add (provision) a new server | Yes |
| POST | `/servers/:serverId/power` | Control server power (on/off) | Yes |

### Example: Login

```bash
curl localhost:8080/auth/login \
  -H "content-type: application/json" \
  -d '{"email": "john.smith@gmail.com", "password": "not-so-secure-password"}'
```
หรือจะใช้การ login ผ่าน browser ที่ port 3000 ก็ได้

### Example: Get all servers

```bash
curl localhost:8080/servers -H "content-type: application/json" -X GET -b "access_token=<JWT_TOKEN>"

```

### Example: Add server

```bash
curl localhost:8080/servers \
  -H "content-type: application/json" \
  -X POST \
  -d '{"sku": "C1-R1GB-D20GB"}' \
  -b "access_token=<JWT_TOKEN>"

```

### Example: Power control

```bash
  curl localhost:8080/servers/{add server id from get api}/power \
  -H "content-type: application/json" \
  -X POST \
  -d '{"action": "off"}' \
  -b "access_token=<JWT_TOKEN>"

```

## Project Structure

```
Cloud-Server-Management/
├── appUtility/config/       # Application configuration
├── frontend/                # Frontend static server (port 3000)
├── services/cloudMgmt/
│   ├── main.go              # Application entrypoint
│   ├── behavior/            # Business logic layer
│   ├── handler/             # HTTP handlers (Fiber)
│   │   └── middleware/      # Auth, CORS, Audit middleware
│   ├── model/               # Data models
│   └── repo/
│       ├── externalRepo/    # External response structs
│       └── gatewayRepo/     # Infra API request/response structs
├── go.mod
└── README.md
```

---


# Data Models

## UserModel

ใช้เก็บข้อมูล user สำหรับ authentication

| Field | Type | Description |
|-------|------|-------------|
| UserId | `string` | User ID (unique identifier) |
| Email | `string` | อีเมลของ user ใช้สำหรับ login |
| PasswordHash | `string` | รหัสผ่านที่ถูก hash ด้วย bcrypt |

**Source:** `services/cloudMgmt/model/user.go`

---

## ServerModel

ใช้เก็บข้อมูล server ที่ถูกสร้าง (provisioned) ในระบบ

| Field | Type | Description |
|-------|------|-------------|
| ServerId | `string` | Server ID (UUID) ใช้อ้างอิงภายในระบบ |
| InfraId | `string` | Infrastructure ID จาก external Infra API |
| Sku | `string` | SKU ของ server เช่น `C1-R1GB-D40GB` |
| IsPowerOn | `bool` | สถานะเปิด/ปิดเครื่อง (`true` = เปิด, `false` = ปิด) |

**Source:** `services/cloudMgmt/model/server.go`

---

## AuditTrailModel

ใช้เก็บ log การกระทำของ user ในระบบ เพื่อ audit/tracking

| Field | Type | Description |
|-------|------|-------------|
| UserID | `int64` | ID ของ user ที่ทำ action |
| ChronoSequence | `string` | ลำดับเวลาของ event |
| Action | `string` | HTTP method ที่ใช้ เช่น `GET`, `POST`, `PUT` |
| ServerID | `string` | Server ID ที่เกี่ยวข้อง (ถ้ามี) |
| Path | `string` | URL path ที่ถูกเรียก |
| IPAddress | `string` | IP address ของ client |
| ResStatus | `int` | HTTP response status code |
| OldValue | `string` | ค่าก่อนเปลี่ยนแปลง |
| NewValue | `string` | ค่าหลังเปลี่ยนแปลง |
| ActionTime | `time.Time` | เวลาที่เกิด action |

**Source:** `services/cloudMgmt/model/audit.go`



Task 3: 
Questions to consider while implementing:
•	What happens if the external service is down during provisioning?
•	How do you ensure the user isn't stuck waiting indefinitely?
(ขอตอบเป็นภาษาไทยเพื่อความรวดเร็วครับ)
Note: คำตอบนี้นำมาจากการค้น code ใน repo internal ของบริษัทปัจจุบันของผมเพื่อดูทุกความเป็นไปได้ + ประสบการณ์ที่เคยเจอในบริษัทปัจจุบัน
•	What happens if the external service is down during provisioning?
Ans: 
จากคำถามที่ว่า ถ้า service infra เกิด crash หรือ down ณ จังหวะ provision server (โดยสมมติฐานว่าสามารถ provision ได้เพียงแค่ช่องทางเดียวจาก api Cloud Server Management service เส้น POST /servers เท่านั้น) 
1.	infra service ควร return error 500 กลับมา -> Cloud Server Management service ควรจะจับ err code 500 + จับ err desc ที่บ่งบอกเฉพาะเจาะจงว่าเกิดการ down ที่ infra service และตอบกลับ user หน้าบ้านด้วย graceful response เช่น

{
  "success": false,
  "error": "Infrastructure service is temporarily unavailable. Please try again later."
}

โดย code ปัจจุบันยังไม่ได้ implement response with error detail จึงควรทำเพิ่มในอนาคตต่อไป
2.	infra service ไม่ตอบกลับมาเลย หรือ slow response จาก infra service -> แบบนี้ควรกำหนด timeout สำหรับการ call gateway service อย่างเช่นในกรณีนี้คือ infra service เพื่อไม่ให้ user หน้าบ้านรอนานเกินไป โดยปัจจุบันที่ผมคุ้นเคยจะใช้ timeout = 7 second ต่อการ call 1 ครั้ง -> หลังจากนั้นก็ตอบ graceful response

3.	infra data สร้างสำเร็จที่ infra service แต่จังหวะจะตอบกลับ server กลับ crash -> ใช้การ implement เพิ่มแบบข้อ 1 แต่หลังจาก infra service กลับมาแล้วควร implement การ retry เพิ่ม เพื่อเช็คความ integrity ของฝั่ง Cloud Server Management service และฝั่ง infra service หรืออีกทางเลือกคือมี batch สิ้นวันที่ออก report เครื่องที่ provisioning ในวันนั้นๆแบบ daily และ cross-check กับ db ของ Cloud Server Management service ทุกวันเพื่อป้องกันข้อมูลตกหล่น

•	How do you ensure the user isn't stuck waiting indefinitely?
Ans: จากข้อนี้คำตอบจะคล้ายกับข้อ 2 ในคำตอบของคำถามแรก นั่นคือการกำหนด timeout ในการ call gateway service (infra service) -> โดยเพื่อให้ user มั่นใจว่า infra server มีปัญหาจริงๆ ตัว Cloud Server Management service ควรจะมีการ retry ก่อนตามจำนวนครั้งที่อยากให้เป็น แล้วค่อย return graceful response กลับไป

หรืออีก 1 วิธีที่ค่อนข้างน่าสนใจคือการเพิ่ม async process อาจจะเป็นการใช้ publish message queue แทนการยิง rest ของเส้น POST /servers โดยหลังจาก publish เพื่อ request การ provision ก็ให้ save status ของ server_id นั้นๆเป็น “processing” จากนั้นก็ให้ทาง infra service publish ผลลัพธ์กลับมา หรือจะยิง rest กลับมาบอกก็ได้ หรือจะให้ยิงเส้น GET /servers เพื่อเช็คก็ทำได้เช่นกัน
