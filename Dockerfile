FROM golang

ADD . /src/kursobot
RUN apt-get update && apt-get install -y supervisor
RUN cd /src/kursobot && make vendor
RUN cd /src/kursobot && make
RUN cd /src/kursobot && make install

EXPOSE 8443

RUN mkdir /var/log/kursobot/supervisord
RUN addgroup kursobot
RUN useradd -g kursobot kursobot

#ENTRYPOINT /usr/bin/supervisord -c /usr/local/kursobot/kursobotd.ini
CMD ["/usr/bin/supervisord", "-c /usr/local/kursobot/kursobotd.ini"]
CMD ["echo", "OK"]
CMD ["cat", "/var/log/kursobot/supervisord/supervisord.log"]
