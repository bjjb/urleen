FROM golang:alpine AS builder
LABEL maintainer "jj@bjjb.org"
WORKDIR src/urleen
COPY . .
ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -i
FROM scratch
WORKDIR /
COPY www /var/www
COPY --from=builder /go/src/urleen/urleen /bin/urleen
EXPOSE 9000
CMD ["/bin/urleen", "-w", "/var/www", "-b", ":9000"]
