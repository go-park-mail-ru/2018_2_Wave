#!/bin/bash

if [ "$#" -eq "0" ]; then
    echo 'No operation selected....'
    echo '   init     -- init project and generate all files'
    echo '-b server   -- generate and build server'
    echo '-b walhalla -- generate and build walhalla'
    echo '   walhalla -- run \ generate+run walhalla'
    echo '-b configs  -- generate and build config generator'
    echo '   configs  -- run \ generate+run config generator'
    echo 'NOTES:'
    echo 'Rebuild and run:'
    echo '-x <cmd> <args> == -b <command> && <cmd> <args>'
    exit 1
fi

# ---------------// ------------ \\---------------
# ---------------||     main     ||---------------
# ---------------\\ ------------ //---------------

conf_dir=./utiles/configs/
conf_bin=./.build/configs

walh_dir=./utiles/walhalla/
walh_bin=./.build/walhalla

-b_walhalla() {
    gen_file utiles/walhalla types
    go_build $walh_dir $walh_bin
}

-b_configs() {
    gen_file utiles configs
    go_build $conf_dir $conf_bin
}

configs() {
    _call_if_missed $conf_bin -b_configs
    _run $conf_bin
}

walhalla() {
    _call_if_missed $walh_bin -b_walhalla
    _run $walh_bin ${@:1:99}
}

-b_server() {
    walhalla server/types/ server/api/
    go_build . ./.build/server
}

init() {
    configs
    walhalla
    -b_server
}

# ---------------// ------------ \\---------------
# ---------------||    Helpers   ||---------------
# ---------------\\ ------------ //---------------
null='/dev/null'

# ubuntu in windows
if grep -q Microsoft /proc/version; then
  ext='.exe'
else
  ext=''
fi

_call_if_missed() {
    if [[ ! -f $1 ]]; then
        $2
    fi
}

_run() {
    $1$ext ${@:2:99}
}

gen_pkg() {
    pushd $1 > $null
    _run easyjson -output_filename $2.gen.go -pkg .
    popd > $null
}

gen_file() {
    pushd $1 > $null
    _run easyjson -output_filename $2.gen.go $2.go
    popd > $null
}

go_build() {
    _run go build -o $2$ext $1
}

go_run() {
    _run go run $1
}

# ---------------// ------------ \\---------------
# ---------------||    router    ||---------------
# ---------------\\ ------------ //---------------

# execute the comand
# <command> <[]args>
if [[ "$1" != -* ]]; then
    $1 ${@:2:99}
# rebuild and execute
# -x <command> <[]args>
elif [[ "$1" = "-x" ]]; then
    -b_$2
    $2 ${@:3:99}
# execute the complex comand
# -<flag> <command> <[]args>
else
    $1_$2 ${@:3:99}
fi
