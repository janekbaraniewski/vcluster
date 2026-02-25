# AGENTS.md

## Cursor Cloud specific instructions

### Project overview

vCluster is an open-source tool that creates virtual Kubernetes clusters inside host cluster namespaces. The repo contains two main binaries:
- `cmd/vcluster` — the syncer/control plane that runs inside the virtual cluster pod
- `cmd/vclusterctl` — the CLI tool for creating and managing virtual clusters

### Tech stack

- **Language**: Go 1.25 (vendored dependencies in `vendor/`)
- **Build**: GoReleaser, Just (task runner)
- **Kubernetes tools**: Helm v3, kubectl, Kind (for local dev clusters)
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

### Docker-in-Docker (Cursor Cloud caveats)

The Cursor Cloud VM runs inside a Firecracker VM with cgroup v1 and limited kernel module support. This affects Kubernetes-in-Docker:

- **Docker storage**: Must use `fuse-overlayfs` (kernel doesn't support overlay2). Daemon config at `/etc/docker/daemon.json`.
- **ip6tables**: The kernel lacks the `ip6table_raw` module. Docker daemon must be started with `"ip6tables": false` in daemon.json to avoid `ip6tables raw table` errors.
- **Kind**: Cannot reliably create clusters because the kubelet inside Kind nodes fails to start (containerd snapshotter / networking issues).
- **k3d (k3s-in-Docker)**: Works if you pass `--k3s-arg "--snapshotter=native@server:0"` and `--k3s-arg "--flannel-backend=host-gw@server:0"`. Without these, k3s fails on overlayfs and flannel vxlan. Pod-to-service networking may still be unreliable.
- **Startup sequence**: After installing Docker, start the daemon with `sudo dockerd &>/tmp/dockerd.log &`, then `sudo chmod 666 /var/run/docker.sock`.
- **E2E tests** (`just e2e`) require a fully functional Kind/k3d cluster, which is unreliable in this environment. Focus on unit tests, chart tests, and build verification.

### Go build notes

- Dependencies are vendored — always use `-mod vendor` flag.
- First build takes ~3 minutes; subsequent builds are cached and much faster.
- GoReleaser snapshot builds skip chart embedding and asset generation hooks.
