FROM alpine:3.7
ADD build/bin/gmail-sender /usr/local/bin/gmail-sender
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ENTRYPOINT ["/usr/local/bin/gmail-sender", "start"]
EXPOSE 80
