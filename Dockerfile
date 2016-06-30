FROM jpetazzo/dind:latest
MAINTAINER Zhao Shuailong <shuailong@tenxcloud.com>

COPY image-sync /usr/local/bin/image-sync
COPY login.sh /usr/local/bin/login.sh
COPY run.sh /usr/local/bin/run.sh

RUN apt-get update -y && apt-get install -y expect && chmod +x /usr/local/bin/image-sync

ENV USERNAME **LinkMe**
ENV PASSWORD **LinkMe**
ENV SRC_REPOSITORY **LinkMe**
ENV DST_REPOSITORY **LinkMe**
ENV SKIP_PUSHED false

CMD ["/bin/bash", "/usr/local/bin/run.sh"]