# c8y-session-bitwarden

A [go-c8y-cli](https://goc8ycli.netlify.app) extension which selects a session from your Bitwarden vault.

The extension provides the `bitwarden` session provider, which `set-session` and `c8y sessions login` use to show a list of the Cumulocity sessions stored in Bitwarden. The selected session is passed to go-c8y-cli, including the password and the TOTP code (if the item has a TOTP secret).

## Install

Requires go-c8y-cli >= [2.57.0](https://github.com/reubenmiller/go-c8y-cli/releases/tag/v2.57.0).

```sh
c8y extensions install reubenmiller/c8y-session-bitwarden
```

c8y downloads the binary for your platform from the latest GitHub release and verifies its checksum. Update it with:

```sh
c8y extensions update session-bitwarden
```

## Setting up Bitwarden

1. Install [bitwarden-cli](https://bitwarden.com/help/cli/)

    **macOS homebrew**

    ```sh
    brew install bitwarden-cli
    ```

2. Log into bitwarden

    ```sh
    bw login
    ```

3. Set the bitwarden session environment variable

    ```sh
    export BW_SESSION="<<token>>"
    ```

    If you're using touchie, then you can save the bitwarden session token using:

    ```sh
    touchie set BW_SESSION
    ```

4. Sync your vault

    ```sh
    bw sync
    ```

## Using it with go-c8y-cli

1. Select a session from the `c8y` folder of your vault

    ```sh
    set-session --from-extension bitwarden
    ```

    Arguments after `--` are passed to the provider, for example to use another folder or to filter the sessions:

    ```sh
    set-session --from-extension bitwarden -- --folder c8y example.com dev
    ```

    The provider needs the `BW_SESSION` environment variable. If it is not set, then go-c8y-cli reads it with your pin entry command. If you're using touchie, then set the pin entry command so the value is read from the macOS keychain using TouchID:

    ```sh
    eval "$(c8y settings update pinEntry "touchie get" --shell auto)"
    ```

2. Optional: Use the provider by default, so `set-session` (without arguments) uses Bitwarden

    ```sh
    eval "$(c8y settings update --shell auto session.provider.type extension)"
    eval "$(c8y settings update --shell auto session.provider.extension bitwarden)"
    ```

## Using it with the go-c8y-cli MCP server

`c8y mcp serve` can show your Bitwarden sessions in its browser session picker, next to your session files. Selecting a session logs in to it, and when the token expires, the session (and a new TOTP code) is read from Bitwarden again.

Requires a go-c8y-cli version which supports the `list` and `select` commands of session providers.

1. Make sure the MCP server can read `BW_SESSION`. The server has no terminal, so either set the pin entry command (recommended), or pass the variable in the MCP client configuration

    ```sh
    eval "$(c8y settings update pinEntry "touchie get" --shell auto)"
    ```

2. Register the server with the provider, e.g. for Claude Code

    ```sh
    claude mcp add --scope user c8y -- c8y mcp serve --session-provider bitwarden
    ```

    Arguments for the provider's list command can be added, for example to use another folder:

    ```sh
    claude mcp add --scope user c8y -- c8y mcp serve --session-provider "bitwarden --folder c8y"
    ```

    If you set the provider as the default (`session.provider.type extension` and `session.provider.extension bitwarden`), then `--session-provider` is not needed.

You can check the sessions which the picker shows, and log in to one without selecting it interactively:

```sh
c8y sessions list --from-extension bitwarden
eval "$(c8y sessions login --from-extension bitwarden --select bitwarden://<item-id>)"
```

You can also run the binary through c8y, for example to see the help:

```sh
c8y session-bitwarden list --help
```

## Development

Build the binary into the extension folder, then install the folder as a local extension. c8y links to the folder, so you only need to rebuild the binary to see your changes.

```sh
go build -o bin/c8y-session-bitwarden .
c8y extensions install .
c8y extensions validate .
```

Releases are published by GoReleaser when a `v*` tag is pushed.

## Install without the extension (go-c8y-cli < 2.57.0)

```sh
go install github.com/reubenmiller/c8y-session-bitwarden@latest
```

Then use it as an external session provider:

```sh
eval "$(c8y settings update --shell zsh session.provider.type external)"
eval "$(c8y settings update --shell zsh session.provider.command "c8y-session-bitwarden list --folder c8y")"
eval "$(c8y settings update --shell zsh session.provider.secrets BW_SESSION)"
```
