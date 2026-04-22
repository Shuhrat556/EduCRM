# EduCRM API — frontend integratsiyasi (Cursor uchun)

Bu hujjat frontend loyihasida Cursor yoki boshqa agentlarga **backend API bilan qanday ishlashni** ko‘rsatadi. Manba kod: `github.com/educrm/educrm-backend`. To‘liq so‘rov/response sxemalari: **Swagger UI** (`/swagger/index.html`) yoki `docs/swagger.yaml`.

---

## 1. Base URL

| Muhit | Misol |
|--------|--------|
| Lokal | `http://localhost:8080` |
| Docker Compose | `http://localhost:8080` |

Barcha versiyalangan REST yo‘llar: **`/api/v1/...`** ostida.

Frontend `.env` (misollar):

```env
# Vite
VITE_API_BASE_URL=http://localhost:8080

# Next.js
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

Ichki HTTP klientda bazaviy URL: `baseURL = import.meta.env.VITE_API_BASE_URL` (yoki Next `process.env.NEXT_PUBLIC_API_BASE_URL`). So‘rovlar: **`${baseURL}/api/v1/...`**.

---

## 2. JSON javob formati (envelope)

Har bir javob shu yoki shu shaklda:

**Muvaffaqiyat:**

```json
{
  "success": true,
  "data": { }
}
```

**Xato:**

```json
{
  "success": false,
  "error": {
    "code": "unauthorized",
    "message": "Human-readable message",
    "kind": "unauthorized"
  }
}
```

Frontendda **avvalo** `success` ni tekshiring; `data` ni `success === true` bo‘lganda oling. Xatolikda `error.kind` HTTP status bilan mos keladi (`validation`, `forbidden`, `not_found`, …).

- `Content-Type: application/json`
- So‘rovlarda ham JSON body yuboriladi (Swaggerda ko‘rsatiladi).

---

## 3. Autentifikatsiya (JWT)

1. **Login:** `POST /api/v1/auth/login`  
   Body: `{ "login": "email@yoki-telefon", "password": "..." }`  
   `login` da `@` bo‘lsa email (katta-kichik harf farqi yo‘q), aks holda telefon.

2. Javob `data` ichida:
   - `access_token` — qisqa muddatli (masalan 15 daq)
   - `refresh_token` — uzun muddatli, serverda hash sifatida saqlanadi
   - `token_type`: `"Bearer"`
   - `expires_in` — soniyada

3. **Himoyalangan endpointlar:** header qo‘shing:
   ```http
   Authorization: Bearer <access_token>
   ```

4. **Access tugaganda:** `POST /api/v1/auth/refresh`  
   Body: `{ "refresh_token": "<refresh_token>" }` — yangi juftlik qaytadi (refresh rotate).

5. **Logout:** `POST /api/v1/auth/logout`  
   - Yoki `Authorization: Bearer <access>` (shu userning barcha refresh sessiyalari o‘chadi)  
   - Yoki body: `{ "refresh_token": "..." }` (bitta sessiya)

6. **Joriy user:** `GET /api/v1/auth/me` — faqat `Authorization: Bearer` bilan.

**Rollar (JWT ichida, `me` da `role` qator sifatida):**  
`super_admin` | `admin` | `teacher` | `student`

**Frontend strategiyasi (tavsiya):** `access_token` ni xotira (memory) yoki `sessionStorage`; `refresh_token` ni `httpOnly` cookie qilish ideal (hozirgi backend oddiy JSON qaytaradi — odatda `localStorage` yoki `sessionStorage` ishlatiladi; XSS xavfini hisobga oling).

---

## 4. CORS va cookie

- Backend CORS `CORS_ALLOWED_ORIGINS` orqali sozlanadi.
- Agar frontend va backend **turli domen**da bo‘lsa, credentialli so‘rovlar uchun backendda `CORS_ALLOW_CREDENTIALS=true` va wildcard `*` emas, aniq origin ro‘yxati kerak.
- Oddiy `Authorization: Bearer` (cookie talab qilmaydigan) rejimda ko‘pincha faqat to‘g‘ri `Origin` ruxsat etilishi kifoya.

---

## 5. Endpointlar xaritasi (`/api/v1`)

Quyidagi jadvalda **yo‘l** va **kim uchun** (ruznomalar) qisqacha berilgan. Batafsil body/query: Swagger.

| Guruh | Metodlar | Yo‘l | Eslatma |
|--------|-----------|------|---------|
| Health | GET | `/health` | Liveness (DBsiz) |
| Health | GET | `/api/v1/health` | Readiness (DB ping) |
| Auth | POST | `/auth/login`, `/auth/refresh`, `/auth/logout` | Login refresh logout |
| Auth | GET | `/auth/me` | Bearer majburiy |
| Users | CRUD | `/users`, `/users/:id`, `/users/:id/status` | `users.manage` (admin/super_admin) |
| Teachers | CRUD | `/teachers`, `/teachers/:id`, `PATCH .../photo` | `teachers.manage` |
| Rooms | CRUD | `/rooms` | `rooms.manage` |
| Groups | CRUD | `/groups` | `groups.manage` |
| Schedules | CRUD | `/schedules` | `schedules.manage` |
| Attendance | POST, GET, PATCH | `/attendance`, `/attendance/:id` | `attendance.manage` |
| Grades | CRUD | `/grades`, `/grades/:id` | `grades.access` |
| Dashboard | GET | `/dashboard/summary` | `dashboard.read` |
| Files | POST, GET, DELETE | `/files`, `/files/register`, `/files/:id` | `files.manage` |
| Notifications | CRUD | `/notifications`, `/notifications/:id`, `PATCH .../read` | inbox/create ruxsatlari alohida |
| AI | POST | `/ai/analytics/...` (5 ta endpoint) | Har biri alohida permission |
| Payments | POST, GET, PATCH, DELETE | `/payments`, `/payments/history`, `/payments/:id` | Staff yoki o‘z tarixini o‘qish |

**Statik fayllar (local storage):**  
URLlar `STORAGE_PUBLIC_BASE_URL` orqali; serverda yo‘l odatda `/static/files/...`.

---

## 6. Swagger va Postman

- Brauzer: **`http://localhost:8080/swagger/index.html`** (`ENABLE_SWAGGER=true` bo‘lishi kerak).
- Postman: Swaggerdan **Import → Link** yoki repodagi `docs/swagger.yaml` ni import qiling.
- “Authorize” uchun Swaggerdagi **Bearer** maydoniga faqat token qiymati (so‘z `Bearer` ni odatda UI o‘zi qo‘shadi).

---

## 7. Xatolik kodlari (qisqa)

| HTTP | `kind` (odatda) | Amal |
|------|------------------|------|
| 400 | validation | Forma xabarlari |
| 401 | unauthorized | Login / refresh |
| 403 | forbidden | Ruxsat / rol |
| 404 | not_found | ID topilmadi |
| 409 | conflict | Dublikat (masalan attendance/grade) |
| 429 | too_many_requests | Rate limit |
| 503 | — | `/api/v1/health` DB o‘chiq |

---

## 8. Cursor uchun qoida (frontend repo)

Frontend loyihasida Cursor **Agents / Rules** ga quyidagiga o‘xshash qo‘shing (yoki `@docs` sifatida bu faylni bog‘lang):

```text
Backend: EduCRM REST API. Base: process.env.VITE_API_BASE_URL or NEXT_PUBLIC_API_BASE_URL + "/api/v1".
Responses: { success, data? } or { success: false, error: { code, message, kind } }.
Auth: POST /auth/login; store access_token; send Authorization: Bearer <token>; refresh via POST /auth/refresh.
Never assume plain JSON without the envelope; read response.data after success===true.
OpenAPI: import from backend docs/swagger.yaml or http://localhost:8080/swagger/doc.json for exact schemas.
```

Agar frontend va backend **bir monorepo**da bo‘lsa, shu faylga `@file` orqali murojaat qiling: `educrm-backend/docs/FRONTEND_API.md`.

---

## 9. Tekshirish ketma-ketligi (minimal)

1. `docker compose up` yoki `go run ./cmd/api` + Postgres.
2. Brauzerda Swagger ochib `POST /auth/login` ni sinash (super_admin yoki boshqa user).
3. Swaggerda **Authorize** ga `access_token` ni yozib xohlagan `GET` ni chaqirish.
4. Frontendda bir xil `baseURL` va `Bearer` bilan `fetch`/`axios` ishlatish.

---

*Batafsil operatsiyalar va DTO maydonlari uchun doim Swagger/OpenAPI manbasiga tayaning.*
 