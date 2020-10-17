#!/bin/bash

# full path to this dir
ROOT_DIR="$( cd "$( dirname "$0" )" && pwd )"

cd "$ROOT_DIR" || exit 1

start_service() {
  kill -9 -q $(ps aux | grep 'go-build' | awk '{print $2}') >/dev/null 2>&1
  go run -race cmd/server/main.go
}

start_watcher() {
  inotifywait -r -m ./ -e close_write,moved_to,create |
  while read -r path action file; do
      if [[ "$file" =~ .*\.go ]]; then # Does the file end with .go?
          start_service &
      fi
  done
}

start_watcher &
start_service &
wait
