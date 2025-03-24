#!/bin/bash
set -e

# 색상 설정
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # 색상 없음

# 루트 권한 확인
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}오류: 이 스크립트는 루트 권한으로 실행해야 합니다.${NC}"
  exit 1
fi

# 변수 설정
INSTALL_DIR="/opt/system-collector"
SERVICE_NAME="system-collector"
SERVICE_USER="system-collector"
SERVICE_GROUP="system-collector"

echo -e "${GREEN}System Collector 제거를 시작합니다...${NC}"

# 서비스 중지 및 비활성화
echo "서비스 중지 및 비활성화 중..."
systemctl stop $SERVICE_NAME || true
systemctl disable $SERVICE_NAME || true
rm -f /etc/systemd/system/system-collector.service
systemctl daemon-reload

# 설치 디렉토리 제거
echo "설치 디렉토리 제거 중..."
rm -rf $INSTALL_DIR

# 사용자 및 그룹 제거
echo "사용자 및 그룹 제거 중..."
if getent passwd $SERVICE_USER > /dev/null; then
    userdel $SERVICE_USER
fi

if getent group $SERVICE_GROUP > /dev/null; then
    groupdel $SERVICE_GROUP
fi

echo -e "${GREEN}제거가 완료되었습니다.${NC}"