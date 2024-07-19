#!/bin/bash

# 设置变量
DOWNLOAD_URL="https://github.com/xibolun/jms_domain_exporter/releases/download/v1.0.0/jms_domain_exporter"
INSTALL_DIR="/opt/jms_domain_exporter"


# 创建安装目录
mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

echo "Downloading ..."
curl -sSL -O "$DOWNLOAD_URL"

# 添加执行权限
chmod 770 jms_domain_exporter

# 获取用户输入
read -p "Enter JMS address: " jms_url
read -p "Enter JMS token: " token
read -p "Enter JMS exporter port: " port
echo

# 配置文件
cat << EOF > config.yml
jms_addr: "$jms_url"
jms_token: "$token"
interval: 10
dial_timeout: 3
http_port: $port
EOF

echo "Installation completed successfully."

$INSTALL_DIR/jms_domain_exporter -c $INSTALL_DIR/config.yml