# Deployment Guide

This guide covers different deployment scenarios for the Excel Template Engine service.

## Table of Contents

- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Production Deployment](#production-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Monitoring and Logging](#monitoring-and-logging)
- [Backup and Recovery](#backup-and-recovery)

## Local Development

### Prerequisites
- Go 1.21+
- MongoDB 7.0+
- Git

### Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/stepanpotapov/Excel-Template-Engine.git
   cd Excel-Template-Engine
   ```

2. **Install dependencies:**
   ```bash
   make install-deps
   ```

3. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

4. **Start MongoDB locally:**
   ```bash
   # macOS
   brew services start mongodb-community

   # Linux
   sudo systemctl start mongod
   ```

5. **Run the service:**
   ```bash
   make run
   ```

6. **Verify:**
   ```bash
   curl http://localhost:8080/health
   ```

## Docker Deployment

### Quick Start

1. **Start all services:**
   ```bash
   docker-compose up -d
   ```

2. **Verify deployment:**
   ```bash
   docker-compose ps
   docker-compose logs app
   ```

3. **Test the service:**
   ```bash
   ./test_workflow.sh
   ```

### Docker Commands

```bash
# Build and start
docker-compose up --build

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f app

# Restart a service
docker-compose restart app

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Custom Docker Build

```bash
# Build image
docker build -t acts-service:v1.0 .

# Run container
docker run -d \
  --name acts-service \
  -p 8080:8080 \
  -e MONGODB_URI=mongodb://mongodb:27017 \
  acts-service:v1.0
```

## Production Deployment

### Environment Setup

1. **Create production environment file:**
   ```bash
   # .env.production
   SERVER_PORT=8080
   SERVER_HOST=0.0.0.0
   BASE_URL=https://acts.yourdomain.com
   
   MONGODB_URI=mongodb://mongodb-user:password@mongodb-host:27017/acts_db?authSource=admin
   MONGODB_DATABASE=acts_db
   MONGODB_COLLECTION=acts
   MONGODB_TIMEOUT=30s
   
   TEMPLATE_PATH=./templates/act_template.xlsx
   GENERATED_PATH=./generated
   
   LOG_LEVEL=info
   LOG_FORMAT=json
   ```

2. **Use production compose file:**
   ```yaml
   # docker-compose.prod.yml
   version: '3.8'
   
   services:
     app:
       image: acts-service:latest
       restart: always
       ports:
         - "8080:8080"
       env_file:
         - .env.production
       volumes:
         - ./generated:/root/generated
         - ./templates:/root/templates
       depends_on:
         - mongodb
       networks:
         - acts-network
       deploy:
         resources:
           limits:
             cpus: '2'
             memory: 1G
           reservations:
             cpus: '1'
             memory: 512M
   
     mongodb:
       image: mongo:7.0
       restart: always
       environment:
         MONGO_INITDB_ROOT_USERNAME: admin
         MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD}
       volumes:
         - mongodb_data:/data/db
         - ./backups:/backups
       networks:
         - acts-network
   
     nginx:
       image: nginx:alpine
       restart: always
       ports:
         - "80:80"
         - "443:443"
       volumes:
         - ./nginx.conf:/etc/nginx/nginx.conf
         - ./ssl:/etc/nginx/ssl
       depends_on:
         - app
       networks:
         - acts-network
   
   volumes:
     mongodb_data:
   
   networks:
     acts-network:
       driver: bridge
   ```

3. **Deploy:**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Nginx Configuration

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream acts_backend {
        server app:8080;
    }

    server {
        listen 80;
        server_name acts.yourdomain.com;
        
        # Redirect to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name acts.yourdomain.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;

        client_max_body_size 10M;

        location / {
            proxy_pass http://acts_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /api/act/download/ {
            proxy_pass http://acts_backend;
            proxy_buffering off;
        }
    }
}
```

## Kubernetes Deployment

### Deployment YAML

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: acts-service
  labels:
    app: acts-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: acts-service
  template:
    metadata:
      labels:
        app: acts-service
    spec:
      containers:
      - name: acts-service
        image: acts-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: acts-secrets
              key: mongodb-uri
        - name: SERVER_PORT
          value: "8080"
        volumeMounts:
        - name: templates
          mountPath: /root/templates
        - name: generated
          mountPath: /root/generated
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: templates
        configMap:
          name: act-templates
      - name: generated
        emptyDir: {}

---
apiVersion: v1
kind: Service
metadata:
  name: acts-service
spec:
  selector:
    app: acts-service
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Deploy to Kubernetes

```bash
# Create namespace
kubectl create namespace acts

# Create secrets
kubectl create secret generic acts-secrets \
  --from-literal=mongodb-uri="mongodb://user:pass@mongo:27017/acts_db" \
  -n acts

# Apply deployment
kubectl apply -f k8s/deployment.yaml -n acts

# Check status
kubectl get pods -n acts
kubectl get services -n acts

# View logs
kubectl logs -f deployment/acts-service -n acts
```

## Monitoring and Logging

### Health Checks

```bash
# Basic health check
curl http://localhost:8080/health

# Response: {"status":"ok","service":"acts-service"}
```

### Log Aggregation

```bash
# Docker logs
docker-compose logs -f app

# Save logs to file
docker-compose logs app > app.log

# Filter errors
docker-compose logs app | grep ERROR
```

### Monitoring with Prometheus (Future)

```yaml
# Prometheus metrics endpoint (to be implemented)
GET /metrics

# Example metrics:
# - acts_created_total
# - acts_generated_total
# - excel_generation_duration_seconds
# - http_requests_total
# - mongodb_operations_total
```

## Backup and Recovery

### MongoDB Backup

```bash
# Create backup
docker-compose exec mongodb mongodump \
  --out /backups/backup-$(date +%Y%m%d-%H%M%S)

# Restore from backup
docker-compose exec mongodb mongorestore \
  /backups/backup-20251104-120000
```

### Automated Backups

```bash
# Add to crontab
0 2 * * * cd /path/to/project && docker-compose exec -T mongodb mongodump --out /backups/backup-$(date +\%Y\%m\%d)
```

### File Backups

```bash
# Backup generated files
tar -czf generated-backup.tar.gz generated/

# Restore
tar -xzf generated-backup.tar.gz
```

## Scaling

### Horizontal Scaling

```bash
# Scale with Docker Compose
docker-compose up -d --scale app=3

# Scale in Kubernetes
kubectl scale deployment acts-service --replicas=5 -n acts
```

### Load Balancing

Use Nginx, HAProxy, or cloud load balancers to distribute traffic across instances.

## Security Checklist

- [ ] Use HTTPS/TLS in production
- [ ] Secure MongoDB with authentication
- [ ] Use environment variables for secrets
- [ ] Enable firewall rules
- [ ] Regular security updates
- [ ] Implement rate limiting
- [ ] Add authentication to API
- [ ] Sanitize file uploads
- [ ] Enable audit logging
- [ ] Use network segmentation

## Performance Optimization

### Database
- Create indexes on frequently queried fields
- Use connection pooling
- Configure appropriate timeouts

### Application
- Enable response compression
- Cache frequently accessed data
- Use worker pools for file generation

### Infrastructure
- Use SSD for file storage
- Adequate RAM for MongoDB
- CDN for static files

## Troubleshooting

### Service won't start

```bash
# Check logs
docker-compose logs app

# Common issues:
# - Port already in use
# - MongoDB not accessible
# - Missing environment variables
```

### MongoDB connection errors

```bash
# Test MongoDB connection
docker-compose exec mongodb mongo --eval "db.adminCommand('ping')"

# Check network
docker-compose exec app ping mongodb
```

### File generation fails

```bash
# Check disk space
df -h

# Check permissions
ls -la generated/

# Verify template exists
ls -la templates/act_template.xlsx
```

## Rollback Procedure

```bash
# Stop current version
docker-compose down

# Restore previous image
docker-compose pull acts-service:previous

# Start with previous version
docker-compose up -d

# Verify
curl http://localhost:8080/health
```

## Update Procedure

```bash
# Pull latest changes
git pull origin main

# Rebuild image
docker-compose build

# Stop old version
docker-compose down

# Start new version
docker-compose up -d

# Verify
./test_workflow.sh
```

---

**For questions or issues, please open an issue on GitHub.**

