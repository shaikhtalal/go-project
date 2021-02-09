#!/bin/bash
while getopts ":m:" opt; do
  case $opt in
    m)
     if go run main.go gql -m $OPTARG; then
       gqlgen generate
       echo "$OPTARG module generated"
      fi
      ;;
  esac
done