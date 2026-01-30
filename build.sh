#!/bin/bash
docker build --platform linux/amd64 -f Dockerfile -t opus-api:latest .
docker save opus-api:latest -o opus-api.tar
echo "Build complete: opus-api.tar"
