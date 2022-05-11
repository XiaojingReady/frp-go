module example

go 1.16

require (
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	server v0.0.0
)

replace server => ./server
