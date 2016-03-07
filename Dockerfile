FROM golang

ADD . /src/kursobot
ADD ./etc/CHECKS /app/CHECKS

RUN apt-get update && apt-get install -y supervisor
RUN cd /src/kursobot && make vendor
RUN cd /src/kursobot && make
RUN cd /src/kursobot && make install

RUN mkdir /var/log/kursobot/supervisord
RUN addgroup kursobot
RUN useradd -g kursobot kursobot

#ENTRYPOINT /usr/bin/supervisord -c /usr/local/kursobot/kursobotd.ini
CMD ["/usr/bin/supervisord", "-c" "/usr/local/kursobot/kursobotd.ini"]

EXPOSE 8443