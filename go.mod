module ire.com/restgwdemo

go 1.17

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.3
	google.golang.org/genproto v0.0.0-20220210181026-6fee9acbd336
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	ire.com/pg v1.1.0
	ire.com/slog v1.3.0
)

require (
	github.com/ChengWu-NJ/dque v1.0.0 // indirect
	github.com/go-pg/zerochecker v0.2.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/bufpool v0.1.11 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace (
	ire.com/pg v1.1.0 => ../pg
	ire.com/slog v1.3.0 => ../slog
)
