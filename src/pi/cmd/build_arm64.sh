#!/bin/bash

go build -o pi

sudo cp pi /usr/local/bin/
sudo mkdir -p  /etc/pi/
sudo cp prod.yml /etc/pi/
sudo cp pi.service /lib/systemd/system/
sudo chmod 644 /lib/systemd/system/pi.service

sudo systemctl daemon-reload
sudo systemctl enable pi.service
sudo systemctl restart  pi.service