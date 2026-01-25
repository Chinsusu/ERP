# Setup Cloudflare Tunnel via Dashboard (Cách Dễ Nhất)

## Bước 1: Truy cập Cloudflare Dashboard

1. Vào https://dash.cloudflare.com
2. Login với account của bạn
3. Chọn domain **xelu.top**

## Bước 2: Tạo Tunnel qua Dashboard

1. Vào menu bên trái: **Zero Trust** → **Access** → **Tunnels**
2. Click **Create a tunnel**
3. Chọn **Cloudflared**
4. Đặt tên tunnel: `erp-tunnel`
5. Click **Save tunnel**

## Bước 3: Install Connector (Copy token)

Dashboard sẽ hiện ra một command với token, dạng:
```bash
cloudflared service install <SUPER_LONG_TOKEN>
```

**COPY token này**, bạn sẽ dùng nó sau.

## Bước 4: Configure Public Hostname

Trong màn hình tunnel configuration:

### Hostname 1: Main ERP
- **Subdomain**: `erp`
- **Domain**: `xelu.top` 
- **Path**: (để trống)
- **Type**: `HTTP`
- **URL**: `nginx:80`

Click **Save hostname**

### Hostname 2: Grafana (Optional)
Click **Add a public hostname**
- **Subdomain**: `grafana.erp`
- **Domain**: `xelu.top`
- **Path**: (để trống)  
- **Type**: `HTTP`
- **URL**: `grafana:3000`

Click **Save hostname**

## Bước 5: Lưu Token

Copy lệnh install từ dashboard, nó sẽ có dạng:
```bash
cloudflared service install eyJhIjoiNzg5Mzg3NDMy...SUPER_LONG_TOKEN
```

Copy **TOÀN BỘ TOKEN** (phần sau `install`)

## Bước 6: Quay lại Terminal

Chạy lệnh mà dashboard cung cấp:
```bash
sudo cloudflared service install <TOKEN_TỪ_DASHBOARD>
```

## Hoặc: Sử dụng Docker (Recommended cho project này)

Tạo file `/opt/erp/.env.cloudflared`:
```bash
TUNNEL_TOKEN=<TOKEN_TỪ_DASHBOARD>
```

Rồi thêm vào docker-compose:
```yaml
services:
  cloudflared:
    image: cloudflare/cloudflared:latest
    container_name: erp-cloudflared
    restart: always
    command: tunnel --no-autoupdate run --token ${TUNNEL_TOKEN}
    environment:
      - TUNNEL_TOKEN=${TUNNEL_TOKEN}
    networks:
      - erp-network
```

---

## ✅ Verify

Sau khi setup xong:

1. Check tunnel status trong dashboard → Should be **HEALTHY**
2. Test URL:
   - https://erp.xelu.top
   - https://grafana.erp.xelu.top

---

**Lưu ý**: Cách này dễ hơn nhiều so với dùng CLI, và dashboard cung cấp UI trực quan hơn!
