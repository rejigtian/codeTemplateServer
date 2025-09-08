#!/bin/bash

# 停止服务
if pgrep -f "template-server" > /dev/null; then
    echo "Stopping template server..."
    pkill -f "template-server"
    sleep 2
    if pgrep -f "template-server" > /dev/null; then
        echo "Failed to stop server gracefully, forcing stop..."
        pkill -9 -f "template-server"
    fi
    echo "Server stopped."
else
    echo "Server is not running."
fi
