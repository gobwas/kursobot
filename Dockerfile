FROM golang
ADD . /src/kursobot
RUN apt-get update && apt-get install -y supervisor
RUN cd /src/kursobot && make vendor && make
RUN ls -la /src/kursobot/bin
COPY /src/kursobot/bin/app /usr/local/kursobot/bin/app
COPY /src/kursobot/etc/config.conf /usr/local/kursobot/kursobot.conf
COPY /src/kursobot/etc/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
#ENTRYPOINT /root/kursobot/bin/app -config=/usr/local/kursobot/kursobot.conf
EXPOSE 8443
CMD ["/usr/bin/supervisord"]