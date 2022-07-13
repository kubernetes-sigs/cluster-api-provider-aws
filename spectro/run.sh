#!/bin/bash

rm ./generated/*

kustomize build --load_restrictor none global > ./generated/core-global.yaml
kustomize build --load_restrictor none base > ./generated/core-base.yaml