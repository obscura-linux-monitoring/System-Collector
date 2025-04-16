FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# 필요한 패키지 설치
RUN apk add --no-cache git

# 의존성 파일 복사 및 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/server ./cmd/server/main.go

# 실행 이미지 생성
FROM alpine:3.16

WORKDIR /app

# 타임존 설정
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul

# 빌드된 바이너리 복사
COPY --from=builder /app/bin/server /app/server

# config 파일 복사
COPY --from=builder /app/bin/configs /app/configs

# 로그 디렉토리 생성
RUN mkdir -p /app/logs

# 포트 노출 (config.yaml에서 지정된 포트와 일치해야 함)
EXPOSE 8087

ENV DB_HOST=host.docker.internal
ENV INFLUXDB_URL=http://host.docker.internal:8086

# 애플리케이션 실행
CMD ["/app/server"]