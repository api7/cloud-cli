Installation
============

In this section, you'll learn how to install Cloud CLI on your local machine.

Build From Source
-----------------

You can also build the Cloud CLI from source.

> Make sure you installed the [Go](https://go.dev/) environment (version >= `1.16`).

```shell
git clone https://github.com/api7/cloud-cli.git
cd /path/to/cloud-cli
make build
```

Download by Go Install
----------------------

> Will suffer from permission problem before this project opening source.
> Remove this note after the project opening source.

Alternatively, you can download the Cloud CLI by using the `go install` command.

```shell
# Install at tree head:
go install github.com/api7/cloud-cli@main
```

See [Versions](https://go.dev/ref/mod#versions) and
[Pseudo-versions](https://go.dev/ref/mod#pseudo-versions) for how to format the
version suffixes.
