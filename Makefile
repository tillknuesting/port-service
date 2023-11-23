PROTOS := $(wildcard pkg/api/*.proto)
GEN_DIR := internal/gen/grpc

.PHONY: proto
proto: 
	@protoc --proto_path=./pkg/api --go_out=$(GEN_DIR) --go_opt=paths=source_relative \
        	--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
        	$(PROTOS)