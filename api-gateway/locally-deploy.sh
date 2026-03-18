echo "Building api-gateway"
docker buildx build --platform linux/amd64,linux/arm64 -t nzhussup/go-api-gateway:latest --push .

echo "Deploying to k8s"
kubectl rollout restart deployment api-gateway

# echo "Pulling image"
# docker pull nzhussup/image-service:latest

# echo "Deploying image-service"
# docker stop image-service
# docker rm image-service
# docker run -d -p 8085:8085 --name image-service nzhussup/image-service:latest
