module example

go 1.10

require (
	client v0.0.0
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace client => ./client
