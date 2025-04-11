#!/bin/bash

url="https://www.post.japanpost.jp/zipcode/dl/kogaki/zip/ken_all.zip"
last_modified_path=$(dirname $0)/../var/last-modified.txt
last_modified=$(cat ${last_modified_path})

status_code=`curl $url -I -H "If-Modified-Since: ${last_modified}" -w %{http_code} -o /dev/null -s`
if [ "$status_code" -eq "200" ]; then
    curl $url -I -s | grep -i "^Last-Modified" | cut -f2- -d' ' | tr -d "\r\n" > $last_modified_path
    echo "updated=true" >> $GITHUB_OUTPUT
fi

exit 0
