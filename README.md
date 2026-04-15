# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and alert command:

```bash
portwatch start --interval 60 --alert "notify-send 'Port Change Detected: {{.Port}}'"
```

Define a baseline of expected open ports:

```bash
portwatch baseline --ports 22,80,443
```

Run a one-time snapshot and print current open ports:

```bash
portwatch scan
```

### Example Output

```
[2024-01-15 10:32:01] INFO  Watching ports | baseline: [22, 80, 443]
[2024-01-15 10:33:01] WARN  New port opened: 8080 (process: python3, pid: 4821)
[2024-01-15 10:34:01] WARN  Port closed: 443
```

## Configuration

portwatch can be configured via a YAML file at `~/.portwatch.yaml`:

```yaml
interval: 30
baseline:
  - 22
  - 80
  - 443
alert: "echo 'Change detected: {{.Port}}'"
```

## License

MIT © [yourusername](https://github.com/yourusername)