# Contributing

Please refer to the [license](LICENSE) before making changes.

[[_TOC_]]

## Issues

Issues are tracked in [GitHub](https://github.com/pedregon/scaffold/issues).

## Getting Started

Assuming that you have the [`go.mod`](go.mod) version of [Golang](https://go.dev/doc/install) or greater installed...

1. Install [`golangci-run`](https://golangci-lint.run/usage/install/).
2. Install [`just`](https://just.systems/man/en/chapter_2.html)
3. Install [`markdownlint`](https://github.com/igorshubovych/markdownlint-cli?tab=readme-ov-file#installation)
4. Install [`git-cliff`](https://git-cliff.org/docs/installation/pypi)
5. Install [`pre-commit`](https://pre-commit.com/#install).

### Task Management

The `just` tool is utilized for managing reposiory tasks.

```shell
just -l
```

## Commits

We use [covenvtional commits](https://www.conventionalcommits.org/).
Please see [`cliff.toml`](cliff.toml) for our supported commit styles.
All commits should be [signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits).

## Versioning

Releases. Git tags shall follow [semantic versioning](https://semver.org/).
