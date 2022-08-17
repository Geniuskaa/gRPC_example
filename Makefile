protoc-prepare-%:
	mkdir "pkg/gen/$*/v1"

protoc-build-%:
	protoc -Icmd/proto/v1 -Ipkg/third_party --go_out=pkg/gen/$*/v1	\
		--go_opt paths=source_relative	\
		--go-grpc_out=pkg/gen/$*/v1	\
		--go-grpc_opt paths=source_relative	\
		$*.proto

protoc-gateway-build-%:
	protoc -Icmd/proto/v1 -Ipkg/third_party	\
		--grpc-gateway_out=pkg/gen/$*/v1 \
    	--grpc-gateway_opt logtostderr=true \
    	--grpc-gateway_opt paths=source_relative \
    	--grpc-gateway_opt generate_unbound_methods=true \
		$*.proto

# При сборке proto файлов использовать эту команду, где после "-" идет
# название вашего proto файла, который вы хотите сбилдить
proto-build-%:
	make protoc-prepare-$* \
		protoc-build-$* \
		protoc-gateway-build-$*

# Метод не работает через make, поэтому писать в консоль вручную
proto-clean-%:
	rm -r pkg/gen/$*/v1

go-install:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc