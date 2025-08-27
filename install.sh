#!/bin/bash

SYSD=/etc/systemd/system
BIN=/usr/local/bin

if [ "$EUID" -ne 0 ]; then
	echo "Error: script must be run as root" >&2
	exit 1
fi

cp dsync.service $SYSD
cp dsync-hourly.service $SYSD
cp dsync-hourly.timer $SYSD
cp dsync $BIN

systemctl daemon-reload
systemctl enable dsync.service
systemctl enable --now dsync-hourly.timer
