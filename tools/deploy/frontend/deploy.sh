# shellcheck disable=SC2086
docker pull ${IMAGE_NAME}
docker rm -f frontend
IMAGE_NAME=${IMAGE_NAME} docker compose up -d

#docker run -d -p 80:80 -p 443:443 --name "${CONTAINER_NAME}" --restart=always -v "${MEOWALIEN_PUBLIC_KEY_PATH}":/certs/meowalien.com.crt -v ${MEOWALIEN_PRIVATE_KEY_PATH}:/certs/meowalien.com.key ${IMAGE_NAME}