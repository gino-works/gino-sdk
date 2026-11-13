##@@ Kubernetes

##@ CRD Generation
run-deepcopy-gen: ## Run deepcopy-gen
	go tool controller-gen object \
	paths=./apis/...

run-generate-crds: ## Run controller-gen for crds
	go tool controller-gen crd paths=./apis/... output:crd:dir=./apis/crds/

k8s-generate: run-deepcopy-gen run-generate-crds fmt ## Run all code gen tasks


##@@ Protobufs

##@ Protobufs
proto-generate: ## Run protoc generation bits
	rm plugin/v1/*.pb.go || true
	go tool buf generate


##@@ All

generate: k8s-generate proto-generate ## Run all generators
	$(MAKE) fmt
