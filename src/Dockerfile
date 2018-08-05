FROM alpine

# glibc is not installed on Alpine.  We might have to sym link it depending 
# on where the image was built
# http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Add ca-certificates
RUN apk add --no-cache ca-certificates

COPY ./linux_amd64/app .
CMD ./app
EXPOSE 8080
