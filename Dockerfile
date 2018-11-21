FROM golang:latest as build
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV CGO_ENABLED 0
RUN go build -o main . 

FROM alpine
RUN apk add curl
COPY --from=build /app/main /app/main
ENTRYPOINT ["/app/main"]
