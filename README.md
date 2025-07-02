# c8y-session-bitwarden

Set a go-c8y-cli session from bitwarden.

## Install

```sh
go install github.com/reubenmiller/c8y-session-bitwarden@latest
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

**Note**: You need to install go-c8y-cli >= [2.52.0](https://github.com/reubenmiller/go-c8y-cli/releases/tag/v2.52.0) to use these instructions

1. In your shell profile, e.g. `~/.zshrc`, then create the following shell function which you can use to use the bitwarden login.

    **Option 1: Change the set-session**

    ```sh
    eval "$(c8y settings update --shell zsh session.provider.type external )"
    eval "$(c8y settings update --shell zsh session.provider.command "c8y-session-bitwarden list --folder c8y")"
    eval "$(c8y settings update --shell zsh session.provider.secrets BW_SESSION)"
    ```

    If you're using touchie, then you can set the pin-entry setting to use touchie to get the BW_SESSION value from the macOS keychain and authenticate using TouchID

    ```sh
    eval "$(c8y settings update pinEntry "touchie get" --shell auto)"
    ```

    **Option 2: Add a new set-session helper**

    ```sh
    eval "$(c8y settings update pinEntry "touchie get" --shell auto)"
    set-session-bitwarden() {
        eval "$(c8y sessions login --from-cmd "c8y-session-bitwarden list --folder c8y" --secrets BW_SESSION "$@")"
    }
```

2. Reload your shell

3. Activate the session

    ```sh
    set-session
    ```
