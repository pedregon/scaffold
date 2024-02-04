# Contributing

Please refer to the [license](LICENSE) before making changes.

## Issues

Issues are tracked in [GitHub](https://github.com/pedregon/scaffold/issues).

## Getting Started

Assuming that you have the [`go.mod`](go.mod) version of [Golang](https://go.dev/doc/install) or greater...

1. Install [`golangci-run`](https://golangci-lint.run/usage/install/).
2. Install [`just`](https://just.systems/man/en/chapter_2.html).
3. Install [`markdownlint`](https://github.com/igorshubovych/markdownlint-cli?tab=readme-ov-file#installation).
4. Install [`git-cliff`](https://git-cliff.org/docs/installation/pypi).
5. Install [`pre-commit`](https://pre-commit.com/#install).

## Task Management

The `just` tool is utilized for managing repository tasks.

```shell
just -l
```

## Commits

> [!Important]
> All commits should be [signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)!

We use [covenvtional commits](https://www.conventionalcommits.org/).
Please see [`cliff.toml`](cliff.toml) for our supported commit styles.

## Pull Requests

> [!WARNING]
> All unit tests must be passing before acceptance.

Our goal is *best effort* unit testing.

## Versioning

Git tags shall follow [semantic versioning](https://semver.org/).
Releases are published to [proxy.golang.org](https://proxy.golang.org/).

### Changelog

To bump the version and update the changelog, run the following:

```shell
just bump
```

### Publishing

1. Run `just tidy lint test` to ensure readiness.
2. Set the latest version in [`Justfile`](Justfile) if not already pushed.
3. Next, update the [`CHANGELOG.md`](CHANGELOG.md) via `just bump`.
4. Create a `git` tag: `git tag v1.0.0`