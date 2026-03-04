#!/bin/bash

# Configuration Parameters
ENDPOINT="http://localhost:8000/user/create"
CONCURRENCY=10000
TOTAL_REQUESTS=10000
RESULTS_FILE=$(mktemp) # Secure temporary file allocation

make_request() {
    # Payload Generation
    RAND_NAME=$(tr -dc 'a-zA-Z' < /dev/urandom | head -c 12)
    RAND_EMAIL="${RAND_NAME,,}$(tr -dc '0-9' < /dev/urandom | head -c 4)@autotest.local"
    RAND_PASS=$(tr -dc 'a-zA-Z0-9!@#$%' < /dev/urandom | head -c 16)

    PAYLOAD=$(cat <<EOF
{"name":"$RAND_NAME","email":"$RAND_EMAIL","password":"$RAND_PASS","confirm_password":"$RAND_PASS"}
EOF
)

    # Execute curl and capture http_code and time_total (in seconds)
    METRICS=$(curl -s -o /dev/null -w "%{http_code},%{time_total}" -X POST "$ENDPOINT" \
        -H "Content-Type: application/json" \
        -d "$PAYLOAD")

    # Output CSV format: HTTP_STATUS,RESPONSE_TIME
    echo "$METRICS"
}

export ENDPOINT
export -f make_request

echo "Executing load vector. Awaiting completion..."

# Execute parallel tasks and route all stdout to the temporary sink
seq 1 $TOTAL_REQUESTS | xargs -n 1 -P $CONCURRENCY -I {} bash -c 'make_request' > "$RESULTS_FILE"

# Process metrics using awk for floating-point calculation
awk -F',' '
BEGIN { 
    total_time = 0; 
    total_hits = 0; 
    success_hits = 0; 
    error_hits = 0 
}
{
    total_hits++
    total_time += $2
    
    # Categorize HTTP response codes
    if ($1 ~ /^2/) {
        success_hits++
    } else {
        error_hits++
    }
}
END {
    if (total_hits > 0) {
        avg_time = total_time / total_hits
    } else {
        avg_time = 0
    }
    
    printf "\n--- Execution Metrics ---\n"
    printf "Total Hits              : %d\n", total_hits
    printf "Successful Hits (2xx)   : %d\n", success_hits
    printf "Failed Hits (Non-2xx)   : %d\n", error_hits
    printf "Avg Response Time       : %.4f seconds\n", avg_time
    printf "-------------------------\n"
}' "$RESULTS_FILE"

# Resource cleanup
rm "$RESULTS_FILE"