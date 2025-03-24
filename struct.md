linux-monitor/
├── cmd/                    # 애플리케이션 진입점
│   └── server/             # 메인 서버 애플리케이션
│       └── main.go         # 메인 진입점
├── configs/                # 구성 파일
│   └── config.yaml         # 앱 구성
├── internal/               # 비공개 애플리케이션 코드
│   ├── collector/          # 시스템 데이터 수집기
│   ├── storage/            # InfluxDB 연결 및 저장 로직
│   └── websocket/          # WebSocket 서버 구현
├── pkg/                    # 공개 라이브러리 코드
│   └── models/             # 데이터 모델 정의
├── scripts/                # 배포 또는 설정 스크립트
├── go.mod                  # Go 모듈 정의
└── go.sum                  # 의존성 체크섬