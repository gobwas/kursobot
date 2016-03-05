FROM golang
ADD . /src/kursobot
RUN apt-get update && apt-get install -y supervisor
RUN cd /src/kursobot && make vendor && make && make install
COPY ./src/kursobot/bin/app /usr/local/kursobot/bin/app
COPY ./etc/config.conf /usr/local/kursobot/kursobot.conf
COPY ./etc/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
#ENTRYPOINT /root/kursobot/bin/app -config=/usr/local/kursobot/kursobot.conf
EXPOSE 8443
CMD ["/usr/bin/supervisord"]