#!/bin/bash
# ============================================
# 自动签发 Let's Encrypt 证书 (Cloudflare DNS-01)
# ============================================
# 前置条件:
#   1. 域名 DNS 托管在 Cloudflare
#   2. 复制 cf.env.example → cf.env 并填入 API Token
#
# 用法:
#   source cf.env && ./init-certs.sh
# ============================================
set -euo pipefail

DOMAIN="${DOMAIN:?请设置 DOMAIN 环境变量，例如: export DOMAIN=api.your-domain.com}"
ROOT_DOMAIN=$(echo "$DOMAIN" | awk -F. '{
    if (NF>2) print $(NF-1)"."$NF;
    else print $0
}')
CERT_DIR="${CERT_DIR:-./certs}"
ACME_IMAGE="neilpang/acme.sh:latest"
CF_TOKEN="${CF_Token:-}"
CF_ACCOUNT="${CF_Account_ID:-}"

# ============================================
# 参数校验
# ============================================
if [ -z "$CF_TOKEN" ] || [ "$CF_TOKEN" = "your-cloudflare-api-token-here" ]; then
    echo "❌ 请先设置 Cloudflare API Token:"
    echo ""
    echo "   1. 访问 https://dash.cloudflare.com/profile/api-tokens"
    echo "   2. 创建 API Token → Edit zone DNS 模板"
    echo "   3. Zone Resources: Include → Specific zone → ${DOMAIN}"
    echo "   4. 将 Token 和 Account ID 填入 cf.env"
    echo "   5. 执行: source cf.env && ./init-certs.sh"
    exit 1
fi

if [ -z "$CF_ACCOUNT" ] || [ "$CF_ACCOUNT" = "your-cloudflare-account-id-here" ]; then
    echo "❌ 请设置 CF_Account_ID (Cloudflare 概览页面右下角)"
    exit 1
fi

echo "============================================"
echo "🔐 签发证书: ${DOMAIN} + *.${ROOT_DOMAIN}"
echo "   验证方式: Cloudflare DNS-01 challenge"
echo "============================================"

# ============================================
# 创建证书输出目录
# ============================================
mkdir -p "${CERT_DIR}/live/${DOMAIN}"

# ============================================
# Step 1: 申请证书 (DNS-01)
# ============================================
echo ""
echo "→ [1/3] 申请 Let's Encrypt 证书..."

docker run --rm \
    -e CF_Token="${CF_TOKEN}" \
    -e CF_Account_ID="${CF_ACCOUNT}" \
    -v "$(pwd)/acme_data:/acme.sh" \
    "${ACME_IMAGE}" \
    --issue \
    --dns dns_cf \
    --dnssleep 30 \
    -d "${DOMAIN}" \
    -d "*.${ROOT_DOMAIN}" \
    --keylength ec-256 \
    --server letsencrypt

echo "✅ 证书申请成功"

# ============================================
# Step 2: 安装证书到 nginx 共享目录
# ============================================
echo ""
echo "→ [2/3] 安装证书到共享目录..."

docker run --rm \
    -v "$(pwd)/acme_data:/acme.sh" \
    -v "$(pwd)/certs:/etc/letsencrypt" \
    "${ACME_IMAGE}" \
    --install-cert \
    -d "${DOMAIN}" \
    --ecc \
    --key-file "/etc/letsencrypt/live/${DOMAIN}/privkey.pem" \
    --fullchain-file "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" \
    --reloadcmd "echo 'cert renewed at \$(date)' >> /etc/letsencrypt/renewal.log"

echo "✅ 证书已安装到 ${CERT_DIR}/live/${DOMAIN}/"

# ============================================
# Step 3: 验证证书
# ============================================
echo ""
echo "→ [3/3] 验证证书..."

CERT_FILE="${CERT_DIR}/live/${DOMAIN}/fullchain.pem"
KEY_FILE="${CERT_DIR}/live/${DOMAIN}/privkey.pem"

if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    EXPIRY=$(openssl x509 -enddate -noout -in "$CERT_FILE" 2>/dev/null | cut -d= -f2)
    SUBJECT=$(openssl x509 -subject -noout -in "$CERT_FILE" 2>/dev/null | cut -d= -f2-)
    echo ""
    echo "📜 证书信息:"
    echo "   域名: ${SUBJECT}"
    echo "   到期: ${EXPIRY}"
    echo "   路径: ${CERT_DIR}/live/${DOMAIN}/"
fi

echo ""
echo "============================================"
echo "✅ 证书签发完成!"
echo ""
echo "下一步:"
echo "1. 在 docker-compose.ghcr.yml 中取消 acme 服务的注释"
echo "2. docker compose up -d acme    (启动自动续期守护进程)"
echo "3. docker compose restart nginx (加载新证书)"
echo "============================================"
