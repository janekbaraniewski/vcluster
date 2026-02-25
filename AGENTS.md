# AGENTS.md

## Cursor Cloud specific instructions

### Project overview

vCluster is an open-source tool that creates virtual Kubernetes clusters inside host cluster namespaces. The repo contains two main binaries:
- `cmd/vcluster` — the syncer/control plane that runs inside the virtual cluster pod
- `cmd/vclusterctl` — the CLI tool for creating and managing virtual clusters

### Tech stack

- **Language**: Go 1.25 (vendored dependencies in `vendor/`)
- **Build**: GoReleaser, Just (task runner)
- **Kubernetes tools**: Helm v3, kubectl, Kind/k3d/k3s (for local dev clusters), DevSpace
- **Container**: Docker (for building images)
- **Lint**: golangci-lint (config in `.golangci.yml`)
- **Tests**: Go test, Helm unittest (chart tests in `chart/tests/`)

### Common commands

See `CONTRIBUTING.md` for full development workflow. Key commands:

| Task | Command |
|------|---------|
| List all Just targets | `just` |
| Lint | `golangci-lint run --timeout 10m` |
| Unit tests | `bash hack/test.sh` |
| Targeted tests | `go test -mod vendor ./pkg/...` |
| Chart tests | `helm unittest chart` |
| Build CLI | `go build -mod vendor -o ./dist/vcluster-cli ./cmd/vclusterctl/main.go` |
| Build syncer | `go build -mod vendor ./cmd/vcluster/...` |
| Build CLI (goreleaser) | `just build-cli-snapshot` |
| Build syncer image (goreleaser) | `just build-snapshot` |
| DevSpace deploy | `devspace deploy -n vcluster` |
| DevSpace dev (interactive) | `devspace dev -n vcluster` |

### Running a local Kubernetes cluster (Cursor Cloud)

The Cursor Cloud VM kernel (cgroup v1, no `xt_comment` module, no overlay2) prevents standard Kind/k3d usage. Use native k3s with specific workarounds:

**Start k3s:**
```bash
sudo k3s server \
  --write-kubeconfig-mode=644 \
  --disable=traefik \
  --snapshotter=native \
  --flannel-backend=host-gw \
  --cluster-dns=10.0.0.2 \
  > /tmp/k3s-server.log 2>&1 &
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
```

**Fix networking (required after k3s start):**
The kernel lacks `xt_comment` so kube-proxy cannot create service routing rules. Fix manually:
```bash
# Allow pod traffic forwarding
sudo iptables -P FORWARD ACCEPT
sudo iptables -A FORWARD -s 10.42.0.0/16 -j ACCEPT
sudo iptables -A FORWARD -d 10.42.0.0/16 -j ACCEPT

# DNAT for Kubernetes API service IP
NODE_IP=$(kubectl get node -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
sudo iptables -t nat -A OUTPUT -d 10.43.0.1/32 -p tcp --dport 443 -j DNAT --to-destination ${NODE_IP}:6443
sudo iptables -t nat -A PREROUTING -d 10.43.0.1/32 -p tcp --dport 443 -j DNAT --to-destination ${NODE_IP}:6443
```

**Deploy vCluster locally:**
```bash
# Build and import image
go build -mod vendor -o vcluster ./cmd/vcluster/...
docker build -t vcluster:dev-local -f Dockerfile.release --build-arg TARGETARCH=amd64 --build-arg TARGETOS=linux .
rm ./vcluster
sudo k3s ctr images import - < <(docker save vcluster:dev-local)

# Create values file
cat > /tmp/vcluster-values.yaml << 'EOF'
controlPlane:
  statefulSet:
    imagePullPolicy: IfNotPresent
    image:
      registry: ""
      repository: vcluster
      tag: dev-local
EOF

# Deploy
./dist/vcluster-cli create my-vcluster -n my-vcluster --create-namespace --connect=false --local-chart-dir ./chart/ -f /tmp/vcluster-values.yaml
```

### Docker setup (Cursor Cloud)

Docker requires `fuse-overlayfs` storage driver and `ip6tables: false`:
```bash
sudo mkdir -p /etc/docker
printf '{\n  "storage-driver": "fuse-overlayfs",\n  "ip6tables": false\n}\n' | sudo tee /etc/docker/daemon.json
sudo dockerd > /tmp/dockerd.log 2>&1 &
sudo chmod 666 /var/run/docker.sock
```

### Go build notes

- Dependencies are vendored — always use `-mod vendor` flag.
- First build takes ~3 minutes; subsequent builds use cache.
- GoReleaser snapshot builds skip chart embedding and asset generation hooks.

### Linear integration

The repo has `hack/linear-sync/` — a CI tool that syncs GitHub PR data to Linear issues at release time. It requires `LINEAR_TOKEN` and `GITHUB_TOKEN` environment variables. It is not a general-purpose Linear client.
