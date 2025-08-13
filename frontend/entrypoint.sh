#!/bin/sh
find /usr/share/nginx/html -type f -name "*.js" -exec sed -i "s|BACKEND_URL|$${BACKEND_URL:-http://localhost:8080}|g" {} +

# Start Nginx
exec nginx -g "daemon off;"