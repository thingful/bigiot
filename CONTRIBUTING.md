# How to contribute

Thank you for your interesting in contributing to this library. There are many
ways to contribute, and we appreciate any and all of them.

* [Feature Requests](#feature-requests)
* [Bug Reports](#bug-reports)
* [Building](#building)
* [Pull Requests](#pull-requests)
* [Coding Conventions](#coding-conventions)
* [Testing](#testing)
* [Helpful Links and Information](#helpful-links-and-information)

## Feature Requests

To request a change to the library please [open an
issue](https://github.com/thingful/bigiot/issues/new) on this repository, and
add the label: *feature request*. If the feature request actually relates to
the BIGIoT marketplace, then your request will be redirected to the appropriate
place.

## Bug Reports

Bugs are of course an inevitability in software, and of course we can't fix
what we don't know about, so please let us know about any bugs or possible bugs
you see. If you aren't sure whether the issue is a bug or not, feel free to
report a bug anyway.

Bugs should be reported by [opening an
issue](https://github.com/thingful/bigiot/issues/new) here on this repo, and
adding the label: *bug*.  While it is not compulsory a good bug report would
look like the following template:

    <short summary of the bug>

    I tried this code:

    <code sample that causes the bug>

    I expected to see this happen: <description>

    Instead this happened: <description>

    I was using this version of the library: <tag or git SHA>

All three components are important: what you did, what you expected, and what
happened instead. If you are able to include version information as well, this
will also be useful in helping us track down the issue.

## Building

This is a standard Go library, but does use some features meaning it is not
currently compatible with Go versions older than 1.9.

If you wish to work on the library, you should follow the following procedure:

* Fork the library into your own account using the procedure described
  [here](https://guides.github.com/activities/forking/)
* Use the standard `go get` command to pull down Thingful's original fork of
  the library into the standard location on your $GOPATH
* Add your fork as a remote to the version you've pulled down using `go get`.
  This will look something like:

        $ cd $GOPATH/src/github.com/thingful/bigiot
        $ git remote add <fork-name> git@github.com:<your-github-name>/bigiot.git

* Now you are free to add whatever feature or bug fix you wish to contribute by
  submitting a [pull request](#pull-requests) back to the original Thingful
  fork.

## Pull Requests

Pull requests are the primary mechanism we use in order to incorporate changes
into the library. GitHub itself has some [great
documentation][about-pull-requests] on using the Pull Request feature.  We use
the "fork and pull" model [described here][development-models], where
contributors push changes to their personal fork and create pull requests to
bring those changes into the source repository.

[about-pull-requests]: https://help.github.com/articles/about-pull-requests/
[development-models]: https://help.github.com/articles/about-collaborative-development-models/

Please make pull requests against the master branch, and if your change
involves code changes, please make sure you have added or updated tests
accordingly.

## Coding Conventions

* All code must be formatted using the standard
  [`gofmt`](https://golang.org/cmd/gofmt/) tool.
* All new code must include test coverage where practicable.

## Testing

Test cases for the library are currently written using
[testify](https://github.com/stretchr/testify), which provides a thin wrapper
around the standard Go `testing` module.  New test cases should be written in
the same form. Please see existing test cases for examples.

## Helpful Links and Information

* [BIG IoT](http://big-iot.eu/)
* [BIG IoT Marketplacel](https://market.big-iot.org/)
* [Thingful](https://www.thingful.net)
