FROM golang

ADD . /src/kursobot
COPY ./etc/CHECKS /app/CHECKS

RUN apt-get update && apt-get install -y supervisor
RUN cd /src/kursobot && make vendor
RUN cd /src/kursobot && make
RUN cd /src/kursobot && make install

RUN addgroup kursobot
RUN useradd -g kursobot kursobot
RUN mkdir /var/log/kursobot/supervisord

#ENTRYPOINT /usr/bin/supervisord -c /mnt/kursobot/config/kursobotd.ini
CMD ["cp", "/mnt/kursobot/tls/server.cert", "/app/server.crt"]
CMD ["chown", "kursobot:kursobot", "/app/server.crt"]
CMD ["/usr/bin/supervisord", "-c", "/mnt/kursobot/config/kursobotd.ini"]

EXPOSE 8443