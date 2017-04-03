#!/bin/bash

# npm install -g gulp --registry=https://registry.npm.taobao.org
if [ ! -d "/app/node_modules" ]; then
    npm install -g --registry=https://registry.npm.taobao.org
else
    echo "npm installed."
fi

if [ ! -d "/app/bower_components" ]; then
    bower --allow-root install --registry=https://registry.npm.taobao.org
else
    echo "bower installed."
fi

gulp compose --allow-root
