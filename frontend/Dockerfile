# Use Nginx to serve static files
FROM nginx:alpine

# Remove default welcome page
RUN rm -rf /usr/share/nginx/html/*

# Copy our frontend files
COPY index.html /usr/share/nginx/html/index.html

# Expose port 80
EXPOSE 80