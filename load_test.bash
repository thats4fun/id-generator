#!/bin/bash

NUM_PROCESSES=1000
WORKER_BINARY="./id-generator"

LOG_DIR="./logs"
mkdir -p "$LOG_DIR"

# Function to run a single worker instance
run_worker() {
    local worker_id=$1
    echo "Starting worker $worker_id..."
    $WORKER_BINARY getid > "$LOG_DIR/worker_$worker_id.log" 2>&1 &
}

for i in $(seq 1 $NUM_PROCESSES); do
    run_worker $i
done

wait

echo "All $NUM_PROCESSES worker processes have completed."
