FROM nginx:alpine

RUN rm -rf /etc/nginx/conf.d/*

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

COPY public/ /usr/share/nginx/html/


# Expose port 80
EXPOSE 80

# Execute the entrypoint script
ENTRYPOINT ["/docker-entrypoint.sh"]