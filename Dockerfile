# Stage 1
FROM golang:1-buster AS stage1

WORKDIR /var/build/go
ENV GOBIN=/var/build/bin
ADD ./src/ ./src/
ADD ./vendor/ ./vendor/
ADD ./go.mod ./
ADD ./go.sum ./
RUN go mod vendor
#ADD ./startScript.sh /var/build/bin/
RUN go build -o /var/build/bin/MangaLibrary ./src/cmd/MangaLibrary/
#RUN go build -o /var/build/bin/jobs ./cmd/jobs/
#RUN go build -o /var/build/bin/load_drivers ./cmd/loadDrivers/

FROM debian:buster AS stage2
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update --fix-missing && \
    apt-get install -yqq --no-install-recommends \
        ca-certificates \
        curl \
        tzdata \
        && \
    apt-get autoclean -yqq && \
    apt-get clean -yqq

FROM stage2 AS stage3
COPY --from=stage1 /var/build/bin/* /usr/local/bin/
#RUN /usr/local/bin/jobs &

#ENTRYPOINT ["/usr/local/bin/startScript.sh"]