#!/bin/bash

task_help(){
    # filter functions with "task_" prefix.
    tasks=()
    while IFS='' read -r line; do tasks+=("$line"); done < <(compgen -A "function" | grep "task_" | sed "s/task_//")
    printf -v tasks "\t%s\n" "${tasks[@]}"

    usage="Usage $ProgName <task> [options]
Tasks:
$tasks
"
    printf "$usage"
}

task=$1
case "$task" in
    "" | "-h" | "--help")
        task_help
        ;;
    *)
        shift
        task_"${task}" "$@"
        ;;
esac
