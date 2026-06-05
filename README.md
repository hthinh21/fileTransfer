# FileTransfer

A lightweight file sharing web app. Upload a file, get a 6-digit code, share it — the recipient enters the code to download. Files expire automatically after 10 minutes.

## Features

- Drag & drop or click-to-select file upload
- Upload progress bar
- 6-digit share code generated on upload
- 10-minute expiry countdown
- Max file size: 10 MB
- Files stored on Cloudflare R2
- Codes stored in Redis (auto-expire with TTL)

## Tech Stack

- **Backend**: Go (standard `net/http`)
- **Storage**: Cloudflare R2
- **Cache**: Redis
- **Frontend**: Vanilla JS + Tailwind CSS

## Project Structure

```
.
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── handler/
│   │   ├── upload.go           # POST /upload
│   │   └── download.go         # GET /download
│   ├── storage/
│   │   ├── r2.go               # Cloudflare R2 client
│   │   └── redis.go            # Redis store
│   └── sharecode/
│       └── code.go             # 6-digit code generator
└── web/
    ├── index.html
    └── js/app.js
```

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Cloudflare R2 bucket
- Redis (included via Docker Compose)

## Setup

**1. Clone the repo**

```bash
git clone <repo-url>
cd fileTransfer
```

**2. Configure environment**

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

```env
R2_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
R2_ACCESS_KEY=your_access_key
R2_SECRET_KEY=your_secret_key
R2_BUCKET=your_bucket_name
REDIS_ADDR=redis:6379
```

**3. Run with Docker Compose**

```bash
docker-compose up -d --build
```

App will be available at `http://localhost:8080`.

## Running Locally (without Docker)

Make sure Redis is running locally, then:

```bash
REDIS_ADDR=localhost:6379 go run ./cmd/server/main.go
```

## API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/upload` | Upload a file, returns `{"code": "XXXXXX"}` |
| `GET` | `/download?code=XXXXXX` | Download file by code |

## How It Works

1. User uploads a file → stored in Cloudflare R2 with a UUID key
2. A 6-digit code is generated and saved in Redis with a 10-minute TTL, pointing to the R2 object
3. Recipient enters the code → server looks up Redis → streams file from R2
4. After download, the code is deleted from Redis
