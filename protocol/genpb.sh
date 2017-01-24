protoc  -I external/ external/*.proto -I common/  --go_out=external/ 

protoc -I rpc/ rpc/*.proto -I common/ --go_out=plugins=grpc:rpc
