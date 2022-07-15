#!/bin/sh
# KEN_ALL.CSV  -> kenall.json
# JIGYOSYO.CSV -> jigyosyo.json
# using ken-all (Copyright (c) 2018 Taiji Inoue. source: https://github.com/inouet/ken-all)

if !(type ken-all > /dev/null); then
    go install github.com/inouet/ken-all@latest
fi

ken-all address KEN_ALL.CSV -t json | sed -e "s/$/,/g" | tr -d '\n' | sed -e "s/^/[/" | sed -e "s/,$/]/" > kenall.json

ken-all office JIGYOSYO.CSV -t json | sed -e "s/$/,/g" | tr -d '\n' | sed -e "s/^/[/" | sed -e "s/,$/]/" | jq -c 'map({
    zip: .zip7,
    pref: .pref,
    city: .city,
    town: .town,
    office: .name,
    office_kana: .kana
})' > jigyosyo.json

exit 0
