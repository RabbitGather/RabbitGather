# shellcheck disable=SC2086
docker pull ${IMAGE_NAME}
docker rm -f frontend
IMAGE_NAME=${IMAGE_NAME} docker compose up -d frontend
