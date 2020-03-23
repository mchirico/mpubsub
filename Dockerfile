FROM golang:1.13.7-alpine3.11 AS build

RUN apk add --no-cache git


WORKDIR $GOPATH/src/github.com/mchirico/mpubsub

# Copy the entire project and build it

COPY . $GOPATH/src/github.com/mchirico/mpubsub

COPY ./credentials /credentials

RUN go build -mod=vendor -o /bin/project

# Special files
RUN mkdir -p /credentials



# This results in a single layer image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /bin/project /bin/project
COPY --from=build /credentials /credentials


EXPOSE 3000
ENTRYPOINT ["/bin/project"]
# Args to project
CMD []

