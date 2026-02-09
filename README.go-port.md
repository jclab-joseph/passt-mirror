# passt Go migration (WIP)

이 변경은 `passt`를 Go로 단계적으로 마이그레이션하기 위한 실행 가능한 기반입니다.

## 목표
- 기존 C 모듈 구조(`*.c`)를 최대한 유지하기 위해 `internal/passt/<module>.go` 1:1 파일 매핑 유지
- TAP 성격의 L2 프레임 전달 파이프라인 구성
- TCP/IP 헬퍼는 gVisor 라이브러리 사용
- Windows / Linux / macOS 크로스 컴파일 가능한 구조 제공

## 현재 구현 상태
- 엔트리 포인트: `cmd/passt-go/main.go`
- 설정/로깅/엔진 루프 구현
- L2 핵심 경로:
  - Ethernet 파서/직렬화 (`packet.go`)
  - TAP 추상화 및 메모리 TAP pair (`tap.go`)
  - MAC 학습 테이블 (`flow.go`)
  - L2 포워더 (`fwd.go`)
  - PCAP 기록기 (`pcap.go`)
- L4 경로: TCP/UDP 프록시 서비스 (`tcp.go`, `udp.go`)

## 테스트 이식/추가
- TCP 전달 검증: `TestTCPProxyForwards`
- UDP echo 전달 검증: `TestUDPProxyForwardsEcho`
- 엔진 종료 검증: `TestEngineShutdown`
- Ethernet 프레임 round-trip 검증
- L2 포워더 프레임 전달 검증
- PCAP 파일 기록 검증

## 빌드 / 검증
```bash
go test ./...
GOOS=windows GOARCH=amd64 go build ./cmd/passt-go
GOOS=darwin GOARCH=amd64 go build ./cmd/passt-go
GOOS=linux GOARCH=amd64 go build ./cmd/passt-go
```
