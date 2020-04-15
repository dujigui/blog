FROM alpine
RUN mkdir -p /app
WORKDIR /app
COPY ./ ./
EXPOSE 8080
ENTRYPOINT ./blog
