#!/bin/bash

cd "$(dirname "$0")/../../.." || exit 1

fyne package --os linux --icon assets/Icon.png

mv chat-client.tar.xz tmp/
