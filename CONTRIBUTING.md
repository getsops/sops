# Contributing to SOPS

The SOPS project welcomes contributions from everyone. Here are a few guidelines
and instructions if you are thinking of helping with the development of SOPS.

## Getting started

- Make sure you have Go 1.19 or greater installed. You can find information on
  how to install Go [here](https://go.dev/doc/install)
- Clone the Git repository and switch into SOPS's directory.
- Run the tests with `make test`. They should all pass.
- If you modify documentation (RST or MD files), run `make checkdocs` to run
  [rstcheck](https://pypi.org/project/rstcheck/) and
  [markdownlint](https://github.com/markdownlint/markdownlint). These should also
  pass. If you need help in fixing issues, create a pull request (see below) and
  ask for help.
- Fork the project on GitHub.
- Add your fork to Git's remotes:
   - If you use SSH authentication:
     `git remote add <your username> git@github.com:<your username>/sops.git`.
   - Otherwise: `git remote add <your username> https://github.com/<your username>/sops.git`.
- Make any changes you want to SOPS, commit them, and push them to your fork.
- **Create a pull request against `main`**, and a maintainer will come by and
  review your code. They may ask for some changes, and hopefully your
  contribution will be merged!

## Guidelines

- Unless it's particularly hard, changes that fix a bug should have a regression
  test to make sure that the bug is not introduced again.
- New features and changes to existing features should be documented, and, if
  possible, tested.

## Communication

If you need any help contributing to SOPS, several maintainers are on the
[`#sops-dev` channel](https://cloud-native.slack.com/archives/C059800AJBT) on
the [CNCF Slack](https://slack.cncf.io).
