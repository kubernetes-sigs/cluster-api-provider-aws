#!/bin/bash

rm ./generated/*

kustomize build --load-restrictor LoadRestrictionsNone global > ./generated/core-global.yaml
kustomize build --load-restrictor LoadRestrictionsNone base > ./generated/core-base.yaml
