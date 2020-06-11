FROM registry.bifund.local:5000/business-base-image
MAINTAINER Liping Wan <guyun_hy@163.com>
EXPOSE 80
COPY ./recordcenter /data/bin/
ENV SERVICENAME="recordcenter"
ENTRYPOINT /data/bin/recordcenter
