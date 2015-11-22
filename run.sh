#!/bin/sh

CompileDaemon -directory="./src/cmd/golang-demo-api/" -build="gb build all" -pattern="(.+\.go)$" -command="bin/golang-demo-api"
