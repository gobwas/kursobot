FROM golang

RUN addgroup kursobot
RUN useradd -g kursobot kursobot
USER kursobot

ADD . /src/kursobot
ADD ./etc/CHECKS /app/CHECKS

RUN apt-get update && apt-get install -y supervisor
RUN cd /src/kursobot && make vendor
RUN cd /src/kursobot && make
RUN cd /src/kursobot && make install

VOLUME /mnt/kursobot/tls
VOLUME /mnt/kursobot/log
VOLUME /mnt/kursobot/config

#ENTRYPOINT /usr/bin/supervisord -c /mnt/kursobot/config/kursobotd.ini
CMD ["/usr/bin/supervisord", "-c", "/mnt/kursobot/config/kursobotd.ini"]

EXPOSE 8443