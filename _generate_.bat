@echo off
FOR %%G IN (api types) DO (
    pushd %%G
    easyjson -pkg .
    popd

    go run ./walhalla %%G
    go fmt ./%%G/
)
