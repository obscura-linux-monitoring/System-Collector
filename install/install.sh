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
CONFIG_DIR="$INSTALL_DIR/configs"
BIN_DIR="$INSTALL_DIR"
SERVICE_NAME="system-collector"
SERVICE_USER="system-collector"
SERVICE_GROUP="system-collector"

echo -e "${GREEN}System Collector 설치를 시작합니다...${NC}"

# 사용자 및 그룹 생성
if ! getent group $SERVICE_GROUP > /dev/null; then
  echo "그룹 $SERVICE_GROUP 생성 중..."
  groupadd -r $SERVICE_GROUP
fi

if ! getent passwd $SERVICE_USER > /dev/null; then
  echo "사용자 $SERVICE_USER 생성 중..."
  useradd -r -g $SERVICE_GROUP -d $INSTALL_DIR -s /sbin/nologin -c "System Collector Service Account" $SERVICE_USER
fi

# 디렉토리 생성
echo "디렉토리 구조 생성 중..."
mkdir -p $BIN_DIR $CONFIG_DIR

# 파일 복사
echo "파일 복사 중..."
cp -f ../bin/server.exec $BIN_DIR/
cp -f ../bin/configs/config.yaml $CONFIG_DIR/
cp -f uninstall.sh $INSTALL_DIR/

# 서비스 파일 복사
echo "서비스 파일 설치 중..."
cp -f ./system-collector.service /etc/systemd/system/

# 권한 설정
echo "권한 설정 중..."
chown -R $SERVICE_USER:$SERVICE_GROUP $INSTALL_DIR
chmod -R 755 $BIN_DIR
chmod 644 $CONFIG_DIR/config.yaml
chmod 644 /etc/systemd/system/system-collector.service

# 서비스 등록 및 시작
echo "서비스 등록 및 시작 중..."
systemctl daemon-reload
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME

echo -e "${GREEN}설치 완료! 서비스 상태:${NC}"
systemctl status $SERVICE_NAME --no-pager

echo ""
echo -e "${GREEN}시스템 로그 확인:${NC}"
echo "journalctl -u $SERVICE_NAME -f"