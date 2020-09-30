#!/bin/bash

rm ./generated/*

kustomize build global > ./generated/core-global.yaml
kustomize build base > ./generated/core-base.yaml
