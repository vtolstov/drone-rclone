FROM library/alpine

RUN apk update && apk add ca-certificates
ADD drone-rclone /
ADD rclone /usr/local/bin/

ENTRYPOINT ["/drone-rclone"]
