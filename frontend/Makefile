build_image:
	docker build -t frontend .

run_image:
	docker compose up -d && docker logs -f --tail 100 frontend

build_and_run_image: build_image run_image
