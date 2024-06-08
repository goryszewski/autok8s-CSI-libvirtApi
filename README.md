# autok8s-CSI-libvirtApi
CSI - sandbox

# run
go run . --endpoint unix:///var/lib/autok8s/csi.sock

# build

# list

k get volumeattachments.storage.k8s.io
k get csinodes.storage.k8s.io master01 -o yaml