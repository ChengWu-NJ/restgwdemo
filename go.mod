module ire.com/restgwdemo

go 1.17

require (
	github.com/felixge/httpsnoop v1.0.3
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.3
	google.golang.org/genproto v0.0.0-20220210181026-6fee9acbd336
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	ire.com/slog v1.3.0
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20211109184856-51b60fd695b3 // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace ire.com/slog v1.3.0 => ../slog
