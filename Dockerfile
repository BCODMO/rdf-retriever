# Create an up to date minimal Debian Jessi build
FROM debian:jessie

RUN apt-get update && \
  apt-get upgrade -y && \
  apt-get install -y --no-install-recommends apt-utils && \
  apt-get install -y ca-certificates && \
  apt-get clean -y && \
  apt-get autoclean -y && \
  apt-get autoremove -y && \
  rm -rf /usr/share/locale/* && \
  rm -rf /var/cache/debconf/*-old && \
  rm -rf /var/lib/apt/lists/* && \
  rm -rf /usr/share/doc/*

ADD ./rdfloader /
COPY ./entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]

