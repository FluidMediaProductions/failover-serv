FROM scratch
COPY ./public-html /app/public-html
COPY ./failover-serv /app/

WORKDIR /app
CMD ["./failover-serv"]

