echo -e "\033[1;32mPulling the latest Docker image...\033[0m"
docker pull ghcr.io/hcd233/aris-blog-api:master

echo -e "\033[1;34mStarting up services with docker-compose...\033[0m"
docker compose -f docker/docker-compose-without-middlewares.yml up -d

echo -e "\033[1;31mPruning unused Docker images...\033[0m"
docker image prune -a -f

echo -e "\033[1;33mDisplaying Docker logs for aris-blog-api...\033[0m"
docker logs -f aris-blog-api --details