FROM nginx:alpine

# Install dos2unix to convert Windows line endings
RUN apk add --no-cache dos2unix

RUN rm -rf /etc/nginx/conf.d/*

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY entrypoint.sh /docker-entrypoint.sh
RUN dos2unix /docker-entrypoint.sh && chmod +x /docker-entrypoint.sh

COPY public/ /usr/share/nginx/html/
COPY src/js/ /usr/share/nginx/html/js/

# Expose port 80
EXPOSE 80

# Execute the entrypoint script
ENTRYPOINT ["/docker-entrypoint.sh"]