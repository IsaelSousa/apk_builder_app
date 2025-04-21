# Apk Android Builder

This application makes it easier generate yours apk bundle of your aplication android.

```bash
# To install libs
go mod tidy

# Generate bin folder (to set on environment variables)
go build -o bin/apk-builder .

# To run (on your terminal)
go run main.go

```

## Examples

![Cmd](https://github.com/IsaelSousa/apk_builder_app/blob/main/src/image/command_line.png?raw=true "CommandLine App")

## How to install

1. Run "go build -o bin/apk-builder .";
2. Save the file in the a safe place;
3. Define a folderpath on your environment variables path (only windows);
4. run "apk-builder" on terminal;
5. Great!

Author: IsaelSousa / PajeDeath