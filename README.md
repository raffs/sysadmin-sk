# Sysadmin Sidekick Operator tool

`sysadmin-sk` is your best friend when performing superhero activities to operate your
complex systems and infrastructure.

**Disclaimer**

sysadmin-sk is not ready yet, still under development and testing and it's not ready to production,
so **USE VERY CAREFULLY**. 

## Features

Currently, sysadmin-sk provides the following features:

* **aws-sqs: move** - Migrate all the messages from one SQS queue to another

## Contributing

There are a few ways to contribute with the `sysadmin-sk` project and you are more
than welcome to do it. Actually, please do it, I need it :)

Use the [Github Issues](https://github.com/raffs/sysadmin-sk/issues) pages to open features, issues
or any other comment regarding this project.

Also feel free to dive into the code and contribute [Pull Requests](https://github.com/raffs/sysadmin-sk/pulls)
and start making this tool better.

To try to keep things organized, please visit the [Github Project](https://github.com/raffs/sysadmin-sk/projects/1) page
for information about issues, upcoming features and discussions.

### Developing

In order to start developing use the following instructions:

```
mkdir -p $GOPATH/src/raffs
cd $GOPATH/src/raffs
git clone https://github.com/raffs/sysadmin-sk
cd sysadmin-sk
```

### Building

To build `sysadmin-sk` command line binary:

```sh
cd $GOPATH/src/raffs/sysadmin-sk
bash ../../build/build.sh
```
