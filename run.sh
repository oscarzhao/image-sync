#!/bin/bash

echo "docker login index.tenxcloud.com"
/usr/bin/expect /usr/local/bin/login.sh

/usr/local/bin/image-sync --src=$SRC_REPOSITORY --dst=$DST_REPOSITORY --skip=$SKIP_PUSHED --v=4