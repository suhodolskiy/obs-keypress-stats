# Obs Keypress Stats

![Preview](.github/preview.gif)

## How to use?

1. Download the latest release for your platform & architecture and decompress the file.
2. Run file.
3. Add Overlay to OBS Studio with URL from application.

## Arguments
| Argument          | Default           | Description                                                      |
|:------------------|:------------------|:-----------------------------------------------------------------|
| `--addr`          | `localhost:8088`  | Address for the HTTP server to listen on.                        |
| `--template`      | *(empty)*         | Path to the HTML template file (e.g. `index.html`).              |
| `--initial-count` | `0`               | Initial value of the counter.                                    |
| `--state-file`    | `./count.txt`     | Path to the state file (used to save and restore program state). |