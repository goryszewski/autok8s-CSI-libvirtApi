#!/bin/bash
echo "[CLEAN]"
sh script/clean.sh

echo "[BUILD]"
sh script/build.sh

echo "[DEPLOY]"
kubectl apply -f ./manifest/csi.yaml

kubectl logs -f --all-containers -l app=csi-libvirt