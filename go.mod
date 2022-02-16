module ire.com/restgwdemo

go 1.17

require (
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.3
	google.golang.org/genproto v0.0.0-20220210181026-6fee9acbd336
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	ire.com/slog v1.3.0
)

require (
	github.com/ChengWu-NJ/dque v1.0.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace ire.com/slog v1.3.0 => ../slog
