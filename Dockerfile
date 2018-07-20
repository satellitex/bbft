FROM alpine:3.8

CMD ["/bbft/bin/linux/bbft"]

WORKDIR /bbft
COPY ./bin ./bin