#!/usr/bin/env bash

res="$(active $@)"
ret_code=$?

if [ "$ret_code" -ne 0 ]; then
    echo "$res"
    return "$ret_code"
fi

if [ -d "$res" ]; then
    cd "$res"
elif [ -f "$res" ]; then
    nvim "$res"
else
    echo "$res"
fi

