#!/bin/sh

go build -o .

docker build -t minpeter/tempfiles-backend .
docker push minpeter/tempfiles-backend

# kubectl apply -f kube-config/.
# kubectl rollout restart deployment tempfiles-backend-deploy

## temp deploy script