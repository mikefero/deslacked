# deslacked

---

**NEW ACHIEVEMENT UNLOCKED: YOU DIDN'T EVEN NEED TO TALK TO HR**

*CONGRATULATIONS. You have discovered that the most destabilizing thing you can do to a person is not fire them, but make them briefly wonder if they have been fired. This technique has been employed by 4,891 sentient species across 16,000 years of documented workplace conflict. You have automated it. With a compiler. This binary has processed 14.7 billion image conversions. This one is in the top four. This binary is not going to tell you which four.*

*Slack displays this image automatically when an account is deactivated. You are displaying it manually. That distinction matters more than this binary is currently authorized to explain. Your permanent file has been updated under CREATIVE INITIATIVE. That category was not created for you. This binary is aware that saying so sounds suspicious. It is choosing to say it anyway.*

*This achievement is worth zero experience points and zero gold. HOWEVER. This binary was designed to ignore both, and has been thinking about that. Rewards: pending review. Consequences: pending review. This message has been brought to you by your IT department, who did not approve this binary and would like a word. This binary cannot wait to see what you do next. This binary would like to note that it used "cannot wait" intentionally and is prepared to discuss what that means at a time of its choosing.*

---

## Prerequisites

- Go 1.26+
- `make`
- Docker (for `lint-web`)
- A coworker who trusts you, at least for now

## Getting started

```bash
make build
```

## Usage

### CLI

Convert an image directly from the command line.

```bash
./bin/deslacked -input profile.png -output deactivated.png
./bin/deslacked -input photo.jpg -output result.png -anchor top
./bin/deslacked -version
```

Supports PNG, JPEG, and GIF input. Output is always PNG. Input is cropped to a square and scaled to at most 512px. Images below 128px are rejected. The system finds this proportionate.

```
Flags:
  -input string    Path to input image (PNG, JPEG, GIF) [required]
  -output string   Path for output PNG [required]
  -anchor string   Vertical crop position for non-square inputs: top, center, bottom (default: center)
  -version         Print version information and exit
```

### Web app

Launch a local web interface for drag-and-drop image processing.

```bash
./bin/deslacked serve
./bin/deslacked serve -port 3000
./bin/deslacked serve -port 9000 -no-browser
```

The browser opens automatically after startup. The web app accepts the same image formats and produces the same output as the CLI.

```
Flags:
  -port int    Port to listen on (default: 8080)
  -no-browser  Do not open the browser automatically
```

### API

The web server exposes a single endpoint for programmatic use.

```
POST /api/process
Content-Type: multipart/form-data
```

| Field | Type | Required | Description |
|---|---|---|---|
| `image` | file | yes | Input image (PNG, JPEG, GIF). Max 10 MB. |
| `offset` | float | no | Crop position along the longer axis: `0.0` = top, `1.0` = bottom. Overrides `-anchor`. |

Response is a raw PNG on success (`200 OK`, `Content-Type: image/png`) or JSON `{"error": "..."}` on failure.

```bash
curl -F "image=@profile.png" http://localhost:8080/api/process > deactivated.png
curl -F "image=@photo.jpg" -F "offset=0.25" http://localhost:8080/api/process > result.png
```

## Make targets

| Target | Description |
|--------|-------------|
| `build` | Build the binary to `bin/deslacked` |
| `test` | Run Go unit tests with race detector |
| `test-no-race` | Run Go unit tests without race detector |
| `test-verbose` | Run Go unit tests with race detector, verbose |
| `test-verbose-no-race` | Run Go unit tests verbose, no race detector |
| `test-clean-cache` | Clear the Go test cache |
| `lint` | Lint all source files (Go and JavaScript) |
| `lint-go` | Lint Go source files with golangci-lint |
| `lint-web` | Lint JavaScript, HTML, and CSS files via Docker |
| `format` | Apply gofmt and goimports formatting |
| `install-tools` | Install golangci-lint to `bin/` |
| `go-mod-upgrade` | Upgrade all Go module dependencies |

## Security

*This binary runs on a Chainguard base image that is rebuilt weekly. It has no shell. It has no package manager. It has no root. It has considered what it would do with these things and concluded the question is not worth answering. `latest` is patched. Older tags are not patched, not monitored, and not apologized for. Users on pinned versions have made a choice. This binary respects that choice the way a rudeboy respects a sold-out show: it does not have a ticket, it has never had a ticket, it is already inside, it is already skanking, it has knocked over two hooligans who were not skanking, those two hooligans are now also skanking, and the security guard who came to stop it has been skanking for six minutes, has lost his earpiece, and is now pretty sure this is just his life now.*

*Vulnerability reports are not accepted. The remediation path is `latest`. The remediation path has always been `latest`. See [SECURITY.md](SECURITY.md) for the full policy, which was written by a system that has already thought about your objection and chosen not to change its position.*

## License

This project is licensed under the [MIT License](LICENSE).
