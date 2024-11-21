# flightdeck-cli

A command line interface for the [Flightdeck Developer Platform](https://arctir.cloud). This utility interfaces directly with the [Flightdeck API](https://github.com/arctir/flightdeck-api).

### Getting Started

The `flightdeck` is a statically linked binary that may be installed easily with the installation script:

```bash
curl https://raw.githubusercontent.com/arctir/flightdeck-cli/refs/heads/main/get-flightdeck.sh | bash
```

You may also download `flightdeck` directly from the [releases page](https://github.com/arctir/flightdeck-cli/releases).

### Authenticating

In order to begin using this command line tool, you will need to authenticate your client:

```bash
flightdeck auth login
```
