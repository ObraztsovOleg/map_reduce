#!/bin/bash

declare -r SERVER_PORT="5000"
declare -r SERVER_ADDRESS="localhost"

# declare -r AGENT_PORTS=(4444 4445 4446 4447)
declare -r GLOW_BASE_DIR=./glow

declare -a pids

declare -r DIRS=("h31" "h55" "h80" "h86")

# Function to be called when user press ctrl-c.
function ctrl_c() {
  for pid in "${pids[@]}"
  do
    echo "Killing ${pid}..."
    kill ${pid}

    # This suppresses the annoying "[pid] Terminated ..." message.
    wait ${pid} &>/dev/null
  done
}

trap ctrl_c SIGINT
echo "You may press ctrl-c to kill all started processes..."

# ./bin/server ${SERVER_PORT} & pids+=($!)
# echo "Started server at ${SERVER_PORT}, pid: $!"

for (( i=0; i < ${#DIRS[@]}; i++ )); do
    ./bin/agent ${DIRS[i]} ${SERVER_PORT} ${SERVER_ADDRESS} & pids+=($!)
    echo "Started agent, pid: $!"
done

echo
echo "Sleep for 10000 seconds..."
sleep 10000