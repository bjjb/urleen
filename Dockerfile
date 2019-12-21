FROM golang AS builder
LABEL maintainer "jj@bjjb.org"
WORKDIR src/urleen
COPY . .
ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -i
FROM scratch
WORKDIR /
COPY www www
COPY --from=builder /go/src/urleen/urleen /urleen
EXPOSE 9000
CMD ["/urleen"]
