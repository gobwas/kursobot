FROM golang
ADD . /src/kursobot
RUN apt-get update && apt-get install -y supervisor
RUN cd /src/kursobot && make vendor
RUN cd /src/kursobot && make
RUN cd /src/kursobot && make install
#ENTRYPOINT /root/kursobot/bin/app -config=/usr/local/kursobot/kursobot.conf
EXPOSE 8443
CMD ["ls -la", "/usr/local/kursobot"]
CMD ["/usr/bin/supervisord", "-c /usr/local/kursobot/kursobotd.conf"]