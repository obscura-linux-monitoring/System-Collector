# System Collector (시스템 콜렉터)

시스템 리소스 데이터를 수집하고 저장하는 서버 프로그램입니다. 클라이언트로부터 WebSocket을 통해 시스템 메트릭스를 수신하여 InfluxDB에 저장합니다.

이 프로그램은 [클라이언트 프로젝트](https://github.com/obscura-linux-monitoring/System-Monitor/releases)의 서버 프로그램입니다.

## 주요 기능

- WebSocket 기반 실시간 데이터 수신
- InfluxDB를 이용한 시계열 데이터 저장
- PostgreSQL을 이용한 명령어 관리
- 다중 클라이언트 동시 연결 지원
- 실시간 시스템 메트릭스 처리
- 확장 가능한 아키텍처

## 시스템 요구사항

- Ubuntu 20.04 LTS 이상
- Go 1.19 이상
- InfluxDB 2.0 이상
- PostgreSQL 12 이상

## 설치 방법

```bash
sudo wget -O /tmp/install.sh https://github.com/obscura-linux-monitoring/System-Collector/releases/download/[version]/install.sh  && chmod +x install.sh && sudo sh /tmp/install.sh [version]
```

## 프로젝트 구조
```
linux-monitor/
├── cmd/ # 애플리케이션 진입점
├── configs/ # 구성 파일
├── internal/ # 비공개 애플리케이션 코드
│ ├── collector/ # 시스템 데이터 수집기
│ ├── storage/ # InfluxDB 연결 및 저장 로직
│ └── websocket/ # WebSocket 서버 구현
├── pkg/ # 공개 라이브러리 코드
│ └── models/ # 데이터 모델 정의
└── scripts/ # 배포 또는 설정 스크립트
```


## 설정 파일

`config.yaml` 파일에서 다음 설정을 구성할 수 있습니다:

- 서버 호스트 및 포트
- InfluxDB 연결 정보
- PostgreSQL 연결 정보
- 데이터 수집 간격

## 개발 환경 설정

1. Go 1.19 이상 설치
2. InfluxDB 2.0 설치 및 구성
3. PostgreSQL 12 설치 및 구성
4. 프로젝트 클론 및 의존성 설치

## 라이선스

MIT License

## 기여하기

1. Fork the Project
2. Create your Feature Branch
3. Commit your Changes
4. Push to the Branch
5. Open a Pull Request

## 문제 해결

자주 발생하는 문제와 해결 방법은 [위키](./wiki/troubleshooting.md)를 참조하세요.

## 연락처

이메일: alkjfgh@gmail.com