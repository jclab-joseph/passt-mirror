# passt Go migration (WIP)

이 변경은 `passt`를 Go로 단계적으로 마이그레이션하기 위한 실행 가능한 기반입니다.

## 목표
- 기존 C 모듈 구조(`*.c`)를 최대한 유지하기 위해 `internal/passt/<module>.go` 1:1 파일 매핑 유지
- TCP/IP 스택은 gVisor netstack 사용 경로 제공 (`-tags gvisor`)
- Windows / Linux / macOS 크로스 컴파일 가능한 구조 제공

## 현재 구현 상태
- 엔트리 포인트: `cmd/passt-go/main.go`
- 설정/로깅/엔진 루프 구현
- TCP/UDP 프록시 서비스 구현 (`tcp.go`, `udp.go`) 및 엔진 연동
- 모듈 매핑 스텁 유지: 아직 포팅되지 않은 C 모듈은 TODO로 추적

## 테스트 이식
기존 테스트 시나리오 중 가장 핵심적인 연결성 검증을 Go 단위 테스트로 이식했습니다.
- TCP 전달 검증: `TestTCPProxyForwards`
- UDP 응답 검증: `TestUDPProxyResponds`

## 빌드 / 검증
```bash
go test ./...
GOOS=windows GOARCH=amd64 go build ./cmd/passt-go
GOOS=darwin GOARCH=amd64 go build ./cmd/passt-go
GOOS=linux GOARCH=amd64 go build ./cmd/passt-go
```

## 알려진 이슈
- `go test -tags gvisor ./...`는 upstream gVisor 모듈의 코드 생성/패키징 이슈로 현재 실패합니다.
