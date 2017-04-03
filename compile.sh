#!/bin/bash

echo "compiling...."
cd /go/src/gateway && go install

cd /go/bin/ && ./gateway