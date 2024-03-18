#!/bin/bash

BMXTIMECODE="./bin/bmxtimecode"
OUTDIR="./testdata"

generate_timecodes() {
    local ratename="$1"
    local fraction="$2"
    local output="$3"

    ${BMXTIMECODE} --output "tc-drop" --rate "${fraction}" all | awk '{$1=$1};1' | sed 's/^[^:]*: //' > "${OUTDIR}/tc-all-${ratename}.txt"
}

generate_timecodes "23_976" "24000/1001"
generate_timecodes "24" "24"
generate_timecodes "29_97" "30000/1001"
generate_timecodes "30" "30"
generate_timecodes "59_94" "60000/1001"
generate_timecodes "60" "60"
