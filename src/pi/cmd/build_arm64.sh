#!/bin/bash

go build -o Pioneer600

sudo cp Pioneer600 /usr/local/bin/
sudo mkdir -p  /etc/Pioneer600/
sudo cp prod.yml /etc/Pioneer600/
sudo cp Pioneer600.service /lib/systemd/system/
sudo chmod 644 /lib/systemd/system/Pioneer600.service

sudo systemctl daemon-reload
sudo systemctl enable Pioneer600.service
sudo systemctl restart  Pioneer600.service