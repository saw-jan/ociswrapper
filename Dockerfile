FROM owncloudci/golang:1.20

WORKDIR /wrapper

COPY . ./

RUN go build -o ./bin/ociswrapper

FROM owncloudci/alpine
COPY --from=0 /wrapper/bin/ociswrapper /usr/bin/ociswrapper

EXPOSE 5000
