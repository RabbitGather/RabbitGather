# shellcheck disable=SC2086
docker pull ${IMAGE_NAME}
docker rm -f "${CONTAINER_NAME}" || true
echo "${MEOWALIEN_PUBLIC_KEY_PATH}"
source ~/meowalin_com_ssh_key_env.sh
docker compose up -d
#docker run -d -p 80:80 -p 443:443 --name "${CONTAINER_NAME}" --restart=always -v "${MEOWALIEN_PUBLIC_KEY_PATH}":/certs/meowalien.com.crt -v ${MEOWALIEN_PRIVATE_KEY_PATH}:/certs/meowalien.com.key ${IMAGE_NAME}
