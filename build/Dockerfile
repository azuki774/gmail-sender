FROM gcr.io/distroless/base-debian10
ADD build/bin/gmail-sender /usr/local/bin/gmail-sender
ENTRYPOINT ["/usr/local/bin/gmail-sender", "start"]
EXPOSE 80
