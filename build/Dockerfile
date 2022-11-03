FROM alpine:3.7
ADD build/bin/gmail-sender /usr/local/bin/gmail-sender
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENTRYPOINT ["/usr/local/bin/gmail-sender", "start"]
EXPOSE 80
