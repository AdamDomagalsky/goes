# Bare metal linux setup

1. hostname setup

sudo hostnamectl set-hostname k8s-control

sudo vi /etc/hosts

172.31.110.168  k8s-control
172.31.98.154   k8s-worker1
172.31.102.145  k8s-worker2

2. Kernel modules

cat << EOF | sudo tee /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF

sudo modprobe overlay &&
sudo modprobe br_netfilter &&
echo 1 | sudo tee /proc/sys/net/ipv4/ip_forward
3. System lvl configuration

cat << EOF| sudo tee /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

4. Reload system
sudo systemctl --system

5. Install containerd

sudo apt-get update && sudo apt-get install -y containerd
6. default config
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
sudo systemctl restart containerd

7. k8s require swapoff to kubelet work properly
sudo swapoff -a

8. [Install k8s](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo deb "https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl
sudo apt-get install -y kubelet kubeadm kubectl



9. Setup kubectl - read kubeadm result
sudo kubeadm init --pod-network-cidr 192.168.0.0/16 --kubernetes-version 1.27.0

10. Calico network plugin
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/master/manifests/calico.yaml

kubeadm token create --print-join-command
sudo kubeadm join 172.31.110.168:6443 --token j8epte.26xughqachazo7w0 --discovery-token-ca-cert-hash sha256:4a081195825a3482b5569d153f84b24e9c9f473cef982527293b1bc1b30d68bb