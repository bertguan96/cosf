if ! command -v protoc-gen-grpc-gateway &> /dev/null; then
    echo "Error: protoc-gen-grpc-gateway not found in PATH." >&2
    exit 1
fi

protoc --proto_path=. --proto_path=./googleapis \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
    cosf/cosf.proto