#!/bin/bash -x

platforms=("darwin/amd64" "linux/amd64" "windows/amd64")
mkdir -p bin

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_os=$GOOS'_'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_os+='.exe'
    fi

	env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -a -ldflags '-w' -o bin/binance_cli_$output_os .
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
