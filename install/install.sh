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

# 버전 매개변수 확인
VERSION=${1:-"dev"}  # 매개변수가 없으면 기본값 "dev" 사용
BASE_URL="https://github.com/obscura-linux-monitoring/System-Collector/releases/download/$VERSION"
DOWNLOAD_EXCUABLE_URL="$BASE_URL/server.exec"
DOWNLOAD_SERVICE_URL="$BASE_URL/system-collector.service"
DOWNLOAD_UNINSTALL_SCRIPT_URL="$BASE_URL/uninstall.sh"
DOWNLOAD_CONFIG_URL="$BASE_URL/config.yaml"

# 변수 설정
INSTALL_DIR="/opt/system-collector"
CONFIG_DIR="$INSTALL_DIR/configs"
BIN_DIR="$INSTALL_DIR"
SERVICE_NAME="system-collector"
SERVICE_USER="system-collector"
SERVICE_GROUP="system-collector"

# 배포판 확인
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VERSION_ID=$(echo $VERSION_ID | tr -d '"')  # 따옴표 제거
    
    # Ubuntu 20.04 버전 체크
    if [ "$OS" != "ubuntu" ] || [ "$VERSION_ID" != "20.04" ]; then
        echo -e "${RED}이 프로그램은 Ubuntu 20.04에서만 실행 가능합니다.${NC}"
        exit 1
    fi
else
    echo -e "${RED}지원되지 않는 리눅스 배포판입니다.${NC}"
    exit 1
fi

# 정리 함수 추가
cleanup() {
    echo -e "${RED}설치 파일 정리 중...${NC}"
    
    # 서비스 파일 제거
    rm -f /etc/systemd/system/$SERVICE_NAME.service 2>/dev/null
    
    # 설치 디렉터리 제거
    rm -rf $INSTALL_DIR 2>/dev/null
    
    echo "정리 완료"
}

# 디렉터리 생성
create_directories() {
    echo "디렉터리 검사 중..."
    if [ -d "$INSTALL_DIR" ]; then
        echo -e "${RED}오류: $INSTALL_DIR 가 이미 존재합니다.${NC}"
        echo "이미 설치되어 있을 수 있습니다. 제거 후 다시 시도해주세요."
        return 1
    fi
    
    echo "디렉터리 생성 중..."
    mkdir -p $BIN_DIR $CONFIG_DIR || return 1
    echo "디렉터리 생성 완료"
    return 0
}

# 파일 다운로드
download_files() {
    echo "실행 파일 다운로드 중..."
    if ! wget -O $BIN_DIR/server.exec $DOWNLOAD_EXCUABLE_URL; then
        echo -e "${RED}실행 파일 다운로드 실패${NC}"
        return 1
    fi
    chmod +x $BIN_DIR/server.exec

    echo "설정 파일 다운로드 중..."
    if ! wget -O $CONFIG_DIR/config.yaml $DOWNLOAD_CONFIG_URL; then
        echo -e "${RED}설정 파일 다운로드 실패${NC}"
        return 1
    fi

    echo "서비스 파일 다운로드 중..."
    if ! wget -O /etc/systemd/system/$SERVICE_NAME.service $DOWNLOAD_SERVICE_URL; then
        echo -e "${RED}서비스 파일 다운로드 실패${NC}"
        return 1
    fi
    
    echo "언인스톨 스크립트 다운로드 중..."
    if ! wget -O $INSTALL_DIR/uninstall.sh $DOWNLOAD_UNINSTALL_SCRIPT_URL; then
        echo -e "${RED}언인스톨 스크립트 다운로드 실패${NC}"
        return 1
    fi
    chmod +x $INSTALL_DIR/uninstall.sh
    
    echo "파일 다운로드 완료"
    return 0
}

# 사용자 및 그룹 생성
create_user_and_group() {
    if ! getent group $SERVICE_GROUP > /dev/null; then
        echo "그룹 $SERVICE_GROUP 생성 중..."
        groupadd -r $SERVICE_GROUP || return 1
    fi

    if ! getent passwd $SERVICE_USER > /dev/null; then
        echo "사용자 $SERVICE_USER 생성 중..."
        useradd -r -g $SERVICE_GROUP -d $INSTALL_DIR -s /sbin/nologin -c "System Collector Service Account" $SERVICE_USER || return 1
    fi
    
    return 0
}

# 권한 설정
set_permissions() {
    echo "권한 설정 중..."
    chown -R $SERVICE_USER:$SERVICE_GROUP $INSTALL_DIR || return 1
    chmod -R 755 $BIN_DIR || return 1
    chmod 644 $CONFIG_DIR/config.yaml || return 1
    chmod 644 /etc/systemd/system/$SERVICE_NAME.service || return 1
    
    return 0
}

# 서비스 등록 및 시작
start_and_enable_service() {
    echo "서비스 등록 및 시작 중..."
    systemctl daemon-reload || return 1
    systemctl enable $SERVICE_NAME || return 1
    systemctl start $SERVICE_NAME || return 1
    
    return 0
}

# 메인 설치 과정
echo -e "${GREEN}System Collector 설치를 시작합니다...${NC}"

create_directories
if [ $? -ne 0 ]; then
    echo -e "${RED}디렉터리 생성 실패. 설치를 중단합니다.${NC}"
    cleanup
    exit 1
fi

create_user_and_group
if [ $? -ne 0 ]; then
    echo -e "${RED}사용자 및 그룹 생성 실패. 설치를 중단합니다.${NC}"
    cleanup
    exit 1
fi

download_files
if [ $? -ne 0 ]; then
    echo -e "${RED}파일 다운로드 실패. 설치를 중단합니다.${NC}"
    cleanup
    exit 1
fi

set_permissions
if [ $? -ne 0 ]; then
    echo -e "${RED}권한 설정 실패. 설치를 중단합니다.${NC}"
    cleanup
    exit 1
fi

start_and_enable_service
if [ $? -ne 0 ]; then
    echo -e "${RED}서비스 시작 및 활성화 실패. 설치를 중단합니다.${NC}"
    cleanup
    exit 1
fi

echo -e "${GREEN}설치 완료! 서비스 상태:${NC}"
systemctl status $SERVICE_NAME --no-pager

echo ""
echo -e "${GREEN}시스템 로그 확인:${NC}"
echo "journalctl -u $SERVICE_NAME -f"