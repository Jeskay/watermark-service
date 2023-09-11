protoc_db_grpc:
	protoc -I api/v1/protos/db api/v1/protos/db/dbsvc.proto --go-grpc_out=api/v1/protos/
protoc_db_msg:
	protoc -I api/v1/protos/db api/v1/protos/db/dbsvc.proto --go_out=api/v1/protos/
protoc_watermark_grpc:
	protoc -I api/v1/protos/watermark api/v1/protos/watermark/watermarksvc.proto --go-grpc_out=api/v1/protos/
protoc_watermark_msg:
	protoc -I api/v1/protos/watermark api/v1/protos/watermark/watermarksvc.proto --go_out=api/v1/protos/
protoc_auth_grpc:
	protoc -I api/v1/protos/auth api/v1/protos/auth/authsvc.proto --go-grpc_out=api/v1/protos/
protoc_auth_msg:
	protoc -I api/v1/protos/auth api/v1/protos/auth/authsvc.proto --go_out=api/v1/protos/