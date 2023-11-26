# Moved: https://github.com/onpremless

```sh
cd cli
go run main.go runtime create ../runtime/golang-1.18/docker/Dockerfile
go run main.go lambda create ../examples/golang-1.18
go run main.go endpoint create
go run main.go lambda start %lambda-name%
curl -d 'test' -X POST http://localhost:8080/%endpoint%
```
