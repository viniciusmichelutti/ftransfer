# ftransfer

Cross-platform LAN file transfer CLI for Mac and Windows. No SMB, no accounts, no servers — auto-discovers peers on your local network via mDNS and streams files over TLS with a live progress bar.

## Usage

On the **receiving** machine:

```bash
ftransfer receive
# Listening as "my-mac" on port 7777 (TLS)
# Output dir: ~/Downloads/ftransfer
# Waiting for transfers... (Ctrl+C to stop)
```

On the **sending** machine — auto-discover and pick a peer:

```bash
ftransfer send ./photos
# Discovering peers on LAN...
#   1) my-mac        192.168.1.10:7777
#   2) work-laptop   192.168.1.22:7777
# Select peer: 1
# Sending 100% |█████████████| (1.2 GB, 24 MB/s)
```

Or skip the picker by targeting directly:

```bash
ftransfer send --to my-mac ./file.txt
ftransfer send --to 192.168.1.10 ./file.txt ./other.zip ./folder/
```

Other commands:

```bash
ftransfer peers      # list LAN peers currently advertising
ftransfer version
ftransfer --help
```

### Flags

`receive`:
- `--out DIR`    output directory (default `~/Downloads/ftransfer`)
- `--port N`     TCP port (default `7777`)
- `--name NAME`  advertised peer name (default: hostname)
- `--yes`        auto-accept all incoming transfers

`send`:
- `--to NAME|IP[:PORT]`  target peer; if omitted, shows interactive picker

## How it works

mDNS advertises `_ftransfer._tcp` on the LAN. The sender dials the receiver over TLS 1.3 (self-signed certs auto-generated on first run, stored in your OS user config dir), exchanges a JSON manifest, waits for the receiver's accept prompt, then streams a tar archive of the requested files/folders.

## Install

Download the binary for your OS/arch from the [latest release](../../releases/latest) and call it on terminal.

Or build from source:

```bash
# Native
go build -o ftransfer ./cmd/ftransfer

# Cross-compile
GOOS=darwin  GOARCH=arm64 go build -o ftransfer      ./cmd/ftransfer
GOOS=windows GOARCH=amd64 go build -o ftransfer.exe  ./cmd/ftransfer
```

