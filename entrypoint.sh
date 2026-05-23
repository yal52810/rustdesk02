#!/bin/sh
# RustDesk API 自动迁移启动脚本
# 容器启动时自动执行数据库迁移

set -e

echo "=========================================="
echo "RustDesk API 启动中..."
echo "=========================================="

# 等待 MySQL 就绪
# 如果配置了 MySQL 则等待并迁移，否则跳过（使用 SQLite）
if [ -n "${DB_HOST}" ]; then
  echo "等待 MySQL 数据库就绪..."
  MAX_TRIES=30
  TRIES=0

  while [ $TRIES -lt $MAX_TRIES ]; do
    if mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -e "SELECT 1" >/dev/null 2>&1; then
      echo "✓ MySQL 已就绪"
      break
    fi
    TRIES=$((TRIES + 1))
    echo "MySQL 未就绪，等待 2 秒... ($TRIES/$MAX_TRIES)"
    sleep 2
  done

  if [ $TRIES -eq $MAX_TRIES ]; then
    echo "✗ MySQL 连接超时"
    exit 1
  fi

  # 执行数据库迁移
  echo ""
  echo "执行数据库迁移..."

  table_exists=$(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -N -e \
    "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='${DB_NAME}' AND TABLE_NAME='user'")

  if [ "$table_exists" = "0" ]; then
    echo "  user 表尚未创建，跳过迁移（将由 API 服务自动创建）"
  else
    add_column_if_not_exists() {
      local column_name=$1
      local column_def=$2
      local exists=$(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -N -e \
        "SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='${DB_NAME}' AND TABLE_NAME='user' AND COLUMN_NAME='${column_name}'")
      if [ "$exists" = "0" ]; then
        echo "  添加列 ${column_name}..."
        mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e \
          "ALTER TABLE \`user\` ADD COLUMN ${column_def}"
      else
        echo "  列 ${column_name} 已存在，跳过"
      fi
    }

    add_index_if_not_exists() {
      local index_name=$1
      local index_def=$2
      local exists=$(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -N -e \
        "SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA='${DB_NAME}' AND TABLE_NAME='user' AND INDEX_NAME='${index_name}'")
      if [ "$exists" = "0" ]; then
        echo "  添加索引 ${index_name}..."
        mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e \
          "CREATE INDEX \`${index_name}\` ON \`user\` ${index_def}"
      else
        echo "  索引 ${index_name} 已存在，跳过"
      fi
    }

    add_column_if_not_exists "relay_server_id" "\`relay_server_id\` bigint unsigned DEFAULT NULL COMMENT '绑定的中继服务器ID'"
    add_index_if_not_exists "idx_relay_server_id" "(\`relay_server_id\`)"
    add_column_if_not_exists "custom_id_server" "\`custom_id_server\` varchar(255) DEFAULT '' COMMENT '自定义ID服务器'"
    add_column_if_not_exists "custom_relay_server" "\`custom_relay_server\` varchar(255) DEFAULT '' COMMENT '自定义中继服务器'"
    add_column_if_not_exists "custom_key" "\`custom_key\` varchar(255) DEFAULT '' COMMENT '自定义服务器公钥'"
  fi

  echo "✓ 数据库迁移完成"
else
  echo "未配置 MySQL (DB_HOST 为空)，使用 SQLite，跳过数据库检查"
fi

echo ""
echo "=========================================="
echo "启动 RustDesk API 服务..."
echo "=========================================="
echo ""

# 启动 API 服务
exec /usr/local/bin/apimain
