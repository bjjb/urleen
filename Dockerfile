FROM golang:alpine
LABEL maintainer "jj@bjjb.org"
ADD . src/github.com/bjjb/urleen
RUN go build github.com/bjjb/urleen
RUN go install github.com/bjjb/urleen
EXPOSE 8089
WORKDIR src/github.com/bjjb/urleen
CMD ["urleen"]
