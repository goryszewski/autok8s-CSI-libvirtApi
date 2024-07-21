# VAR

REPO=repo.internal:5000
IMAGE=csi-libvirt
TAG=latest

# TARGET

build:
	
	DOCKER_BUILDKIT=0 docker build -t $(REPO)/$(IMAGE):$(TAG) .
	docker push $(REPO)/$(IMAGE):$(TAG)
clean:
	kubectl delete -f ./manifest/test-csi.yaml
	kubectl delete -f ./manifest/csi.yaml

deploy:
	kubectl apply -f ./manifest/csi.yaml

re: clean build deploy

status:
	@echo "[volumeattachments.storage.k8s.io]"
	kubectl get volumeattachments.storage.k8s.io
	@echo "[pv]"
	kubectl get pv
	@echo "[PVC]"
	kubectl get pvc
	@echo "[POD]"
	kubectl get pod

log:
	kubectl logs -f --all-containers -l app=csi-libvirt