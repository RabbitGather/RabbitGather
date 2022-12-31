# shellcheck disable=SC2086
docker pull ${IMAGE_NAME}
docker rm -f ${CONTAINER_NAME} || true
IMAGE_NAME=${IMAGE_NAME} CONTAINER_NAME=${CONTAINER_NAME} docker compose up -d deploy

