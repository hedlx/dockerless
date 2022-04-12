FROM golang:1.18-alpine AS bootstrap

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY . .

RUN go mod download
RUN go build -o /lambda

# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

FROM alpine AS lambda

WORKDIR /
COPY --from=bootstrap /lambda /lambda

HEALTHCHECK CMD curl --fail http://localhost:3000/health || exit 1 

ENTRYPOINT ["/lambda"]
