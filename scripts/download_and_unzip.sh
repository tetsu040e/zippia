#!/bin/bash

ken_all_url="https://www.post.japanpost.jp/zipcode/dl/kogaki/zip/ken_all.zip"
jigyosyo_url="https://www.post.japanpost.jp/zipcode/dl/jigyosyo/zip/jigyosyo.zip"

curl $ken_all_url -s -o KEN_ALL.ZIP
unzip KEN_ALL.ZIP KEN_ALL.CSV

curl $jigyosyo_url -s -o JIGYOSYO.ZIP
unzip JIGYOSYO.ZIP JIGYOSYO.CSV

exit 0
