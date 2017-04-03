#!/bin/bash

if [ ! -d "/app/node_modules" ]; then
    npm install
    bower install
else
    echo "npm and bower installed."
fi

gulp compose
