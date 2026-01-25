#!/bin/bash
# server-setup.sh - Production Server Hardening for ERP Cosmetics

# 1. Update system
apt update && apt upgrade -y

# 2. Install required packages
apt install -y \
  docker.io \
  docker-compose-v2 \
  nginx \
  certbot \
  python3-certbot-nginx \
  fail2ban \
  ufw \
  htop \
  ncdu

# 3. Configure firewall (UFW)
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow http
ufw allow https
ufw --force enable

# 4. Configure fail2ban
cat > /etc/fail2ban/jail.local << EOF
[sshd]
enabled = true
maxretry = 3
bantime = 3600

[nginx-http-auth]
enabled = true
EOF

systemctl enable fail2ban
systemctl start fail2ban

# 5. Create app user and directories
useradd -m -s /bin/bash erp
usermod -aG docker erp

mkdir -p /opt/erp
mkdir -p /var/log/erp
mkdir -p /backups

chown -R erp:erp /opt/erp
chown -R erp:erp /var/log/erp
chown -R erp:erp /backups

echo "✓ Production server setup completed!"
