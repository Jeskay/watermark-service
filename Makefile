protoc_grpc:
	protoc -I api/v1/protos api/v1/protos/watermark/watermarksvc.proto api/v1/protos/auth/authsvc.proto api/v1/protos/picture/picturesvc.proto --go-grpc_out=./api/v1/protos --go-grpc_opt=paths=source_relative
protoc_msg:
	protoc -I api/v1/protos api/v1/protos/auth/authsvc.proto api/v1/protos/picture/picturesvc.proto api/v1/protos/watermark/watermarksvc.proto --go_out=./api/v1/protos --go_opt=paths=source_relative