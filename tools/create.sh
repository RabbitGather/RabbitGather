#!/bin/bash


sudo snap install core; sudo snap refresh core
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot certonly --standalone



sudo apt-get update
sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg


echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null


sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin

sudo groupadd docker

sudo usermod -aG docker a_meowalien

newgrp docker


sudo groupadd certbot
sudo usermod -aG certbot a_meowalien
sudo chgrp certbot -R /etc/letsencrypt/live/

sudo usermod -aG certbot docker_registry

sudo useradd -M docker_registry


sayken kingkingjin


gcloud compute ssh --zone asia-east1-b a_meowalien@instance-1 -- '
    cd /home/a_meowalien/rabbit_gather/frountend
    docker pull meowalien.com:5000/frontend:main
    docker stop frontend
    docker rm frontend
    docker run -d -p 80:80 -p 443:443 --name frontend --restart=always \
          -v /etc/letsencrypt/live/meowalien.com/fullchain.pem:/certs/meowalien.com.crt \
          -v /etc/letsencrypt/live/meowalien.com/privkey.pem:/certs/meowalien.com.key \
          meowalien.com:5000/frontend:main

    '