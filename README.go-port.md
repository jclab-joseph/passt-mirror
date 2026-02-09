# passt Go migration (WIP)

이 디렉터리 변경은 `passt` C 구현을 Go로 단계적으로 마이그레이션하기 위한 초기 골격입니다.

## 목표
- 기존 C 모듈 구조(`*.c`)를 최대한 유지하기 위해 `internal/passt/<module>.go` 1:1 파일 매핑을 구성
- TCP/IP 스택은 gVisor netstack (`gvisor.dev/gvisor/pkg/tcpip/...`) 사용
- Windows / Linux / macOS 크로스 컴파일 가능한 순수 Go 중심 구조 제공

## 현재 상태
- 엔트리 포인트: `cmd/passt-go/main.go`
- 설정/로깅/엔진 루프/기본 netstack 초기화 구현
- 나머지 모듈은 C 코드 대응 파일로 스텁 생성 후 TODO 주석으로 포팅 지점 명시

## 빌드
```bash
go build ./cmd/passt-go
```

## 다음 단계
1. `tap.c`, `packet.c`, `flow.c` 순으로 데이터패스 포팅
2. `tcp.c`, `udp.c`, `icmp.c` 포팅 후 통합 테스트 이식
3. `vhost_user.c`, `virtio.c` 포팅으로 qemu 통합 완성
