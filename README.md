# logslice

Fast log file splitter and time-range extractor for large structured log archives.

---

## Installation

```bash
go install github.com/youruser/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/logslice.git && cd logslice && go build ./...
```

---

## Usage

Extract log entries within a specific time range:

```bash
logslice --input app.log --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --output slice.log
```

Split a large log file into chunks by hour:

```bash
logslice split --input app.log --interval 1h --output-dir ./chunks/
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Path to the source log file | required |
| `--from` | Start of time range (RFC3339) | beginning of file |
| `--to` | End of time range (RFC3339) | end of file |
| `--output` | Output file path | stdout |
| `--interval` | Chunk size for split mode | `1h` |
| `--format` | Log timestamp format | `rfc3339` |

### Supported Log Formats

- JSON structured logs
- Common log format (CLF)
- Custom patterns via `--pattern` flag

---

## License

MIT © 2024 youruser