# https://github.com/GoogleContainerTools/distroless
# https://docs.docker.com/language/golang/build-images/

FROM golang:1.18-bullseye as BUILD

WORKDIR /src

COPY ["./", "/src/"]

RUN go mod download \
  && go build -o /src/restapi-go


FROM gcr.io/distroless/base-debian11:nonroot as FINAL

COPY --from=build ["/src/restapi-go", "/usr/local/bin/"]

ENTRYPOINT [ "restapi-go" ]