FROM busybox
RUN mkdir /app
COPY . /app
WORKDIR /app
CMD ["./chat_server", "-app=login", "[login@1]"]