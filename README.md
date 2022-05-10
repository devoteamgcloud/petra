# Petra

### Local Machine

Follow these steps if you are OK installing and using Go on your machine.

1. Install [Go](https://golang.org/doc/install).
1. Install [Visual Studio Code](https://code.visualstudio.com/).
1. Install [Go extension](https://code.visualstudio.com/docs/languages/go).
1. Install [Commitizen](https://github.com/commitizen-tools/commitizen)
1. Install [pre-commit](https://github.com/pre-commit/pre-commit)
1. Clone and open this repository.
1. `F1` -> `Go: Install/Update Tools` -> (select all) -> OK.


## Build

### Terminal

- `make` - execute the build pipeline.
- `make help` - print help for the [Make targets](Makefile).

### Visual Studio Code

`F1` → `Tasks: Run Build Task (Ctrl+Shift+B or ⇧⌘B)` to execute the build pipeline.

## Documentation

### CLI
* [petra cli](cli/doc/petra.md) - CLI to upload / remove / upload a terraform module to a private registry (Google Cloud Storage bucket).
### CLI
* [petra server](server/doc/petra.md) - Server to get a terraform module versions / get a signed URL to download a module from a private registry (Google Cloud Storage bucket).

## Release

The release workflow is triggered each time a tag with `v` prefix is pushed.

_CAUTION_: Make sure to understand the consequences before you bump the major version. More info: [Go Wiki](https://github.com/golang/go/wiki/Modules#releasing-modules-v2-or-higher), [Go Blog](https://blog.golang.org/v2-go-modules).

## Maintainance

Remember to update Go version in [.github/workflows](.github/workflows) and [Makefile](Makefile). 

Notable files:

- [.github/workflows](.github/workflows) - GitHub Actions workflows,
- [.github/dependabot.yml](.github/dependabot.yml) - Dependabot configuration,
- [.golangci.yml](.golangci.yml) - golangci-lint configuration,
- [.goreleaser.yml](.goreleaser.yml) - GoReleaser configuration,
- [Makefile](Makefile) - Make targets used for development, [CI build](.github/workflows) and [.vscode/tasks.json](.vscode/tasks.json),
- [go.mod](go.mod) - [Go module definition](https://github.com/golang/go/wiki/Modules#gomod),
- [tools.go](tools.go) - [build tools](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module).

## FAQ

### Why Visual Studio Code editor configuration

Developers that use Visual Studio Code can take advantage of the editor configuration. While others do not have to care about it. Setting configs for each repo is unnecessary time consuming. VS Code is the most popular Go editor ([survey](https://blog.golang.org/survey2019-results)) and it is officially [supported by the Go team](https://blog.golang.org/vscode-go).

You can always remove the [.devcontainer](.devcontainer) and [.vscode](.vscode) directories if it really does not help you.


### How can I customize the release or add deb/rpm/snap packages, Homebrew Tap, Scoop App Manifest etc

Take a look at GoReleaser [docs](https://goreleaser.com/customization/) as well as [its repo](https://github.com/goreleaser/goreleaser/) how it is dogfooding its functionality.

## Contributing

Simply create an issue or a pull request.