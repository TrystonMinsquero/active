#!/usr/bin/env bash

res="$(active $@)"
ret_code=$?

# should edit cache
if [ "$ret_code" -eq 4 ]; then
    if [[ -v EDITOR ]]; then
        $EDITOR "$res"
    else
        echo "You can set the editor enviornment variable to pick which editor to use"
        vim "$res"
    fi 
    return $?
fi

# should fuzzy find
if [ "$ret_code" -eq 3 ]; then
    if command -v fzf &> /dev/null; then
        fzf_res="$(fzf --query=$@ --select-1 --exit-0 <<< $res)"
        cd "$fzf_res"
    else
        echo "Did you mean to say one of these?"
        echo "$res"
        return 1
    fi
    return $?
fi

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

return $?

