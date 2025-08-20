FROM nginx:alpine

RUN apk add --no-cache dos2unix

RUN rm -rf /etc/nginx/conf.d/*

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY entrypoint.sh /docker-entrypoint.sh
RUN dos2unix /docker-entrypoint.sh && chmod +x /docker-entrypoint.sh

COPY public/ /usr/share/nginx/html/
COPY src/js/ /usr/share/nginx/html/js/

EXPOSE 80

ENTRYPOINT ["/docker-entrypoint.sh"]
