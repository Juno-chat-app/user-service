build:
	protoc -I=. user.proto --go_out=plugins=grpc:.
	protoc -I=. user_message.proto --go_out=plugins=grpc:.
	protoc -I=. user_data.proto --go_out=plugins=grpc:.
	ls *.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'