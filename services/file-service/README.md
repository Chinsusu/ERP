# File Service

File Service quản lý uploads, downloads, và storage cho ERP Cosmetics với MinIO integration.

## Features

- ✅ **File Upload** - Single/multiple file upload
- ✅ **MinIO Storage** - S3-compatible object storage
- ✅ **File Categories** - Configurable validation rules
- ✅ **Presigned URLs** - Secure download links
- ✅ **Entity Attachment** - Link files to entities

## Quick Start

```bash
# Start MinIO (for development)
docker run -d --name minio \
  -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  minio/minio server /data --console-address ":9001"

# Set environment
export DB_HOST=localhost DB_PORT=5432 DB_USER=erp_user \
       DB_PASSWORD=erp_password DB_NAME=erp_file \
       MINIO_ENDPOINT=localhost:9000 \
       MINIO_ACCESS_KEY_ID=minioadmin \
       MINIO_SECRET_ACCESS_KEY=minioadmin

# Run migrations
make migrate-up

# Run service
make run
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/files/upload` | Upload single file |
| POST | `/api/v1/files/upload/multiple` | Upload multiple files |
| GET | `/api/v1/files/categories` | List file categories |
| GET | `/api/v1/files/:id` | Get file metadata |
| GET | `/api/v1/files/:id/download` | Download file |
| GET | `/api/v1/files/:id/url` | Get presigned URL |
| DELETE | `/api/v1/files/:id` | Delete file |
| GET | `/api/v1/files/entity/:type/:id` | Get files by entity |

## File Categories

| Category | Extensions | Max Size |
|----------|------------|----------|
| DOCUMENT | pdf, doc, docx, xls, xlsx | 10 MB |
| IMAGE | jpg, jpeg, png, gif, webp | 5 MB |
| CERTIFICATE | pdf | 5 MB |
| CONTRACT | pdf, doc, docx | 20 MB |
| REPORT | pdf, xlsx, csv | 50 MB |
| AVATAR | jpg, png | 2 MB |
| PRODUCT_IMAGE | jpg, jpeg, png, webp | 5 MB |
| QC_PHOTO | jpg, jpeg, png | 10 MB |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8091 | HTTP port |
| `MINIO_ENDPOINT` | localhost:9000 | MinIO endpoint |
| `MINIO_ACCESS_KEY_ID` | minioadmin | Access key |
| `MINIO_SECRET_ACCESS_KEY` | minioadmin | Secret key |
| `MINIO_USE_SSL` | false | Use SSL |
| `MAX_UPLOAD_SIZE` | 52428800 | Max upload (50MB) |

## License

Proprietary - ERP Cosmetics System
