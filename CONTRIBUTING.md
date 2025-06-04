# Contributing

Wow, we really appreciate that you even looked at this section! We are trying to make the worlds best atomic building blocks for financial services that accelerate innovation in banking and we need your help!

You only have a fresh set of eyes once! The easiest way to contribute is to give feedback on the documentation that you are reading right now. This can be as simple as sending a message to our Google Group with your feedback or updating the markdown in this documentation and issuing a pull request.

Stability is the hallmark of any good software. If you find an edge case that isn't handled please open an GitHub issue with the example data so that we can make our software more robust for everyone. We also welcome pull requests if you want to get your hands dirty.

Have a use case that we don't handle; or handle well! Start the discussion on our Google Group or open a GitHub Issue. We want to make the project meet the needs of the community and keeps you using our code.

Please review our [Code of Conduct](CODE_OF_CONDUCT.md) to ensure you agree with the values of this project.

We use GitHub to manage reviews of pull requests.

* If you have a trivial fix or improvement, go ahead and create a pull request, addressing (with `@...`) one or more of the maintainers (see [AUTHORS.md](AUTHORS.md)) in the description of the pull request.

* If you plan to do something more involved, first propose your ideas in a Github issue. This will avoid unnecessary work and surely give you and us a good deal of inspiration.

* Relevant coding style guidelines are the [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments) and the _Formatting and style_ section of Peter Bourgon's [Go: Best Practices for Production Environments](http://peter.bourgon.org/go-in-production/#formatting-and-style).

* When in doubt follow the [Go Proverbs](https://go-proverbs.github.io/)

* Checkout this [Overview of Go Tooling](https://www.alexedwards.net/blog/an-overview-of-go-tooling) by Alex Edwards

## Getting the code

We recommend using additional git remote's for pushing/pulling code. Go cares about where the `iso8583` project lives relative to `GOPATH`.

To pull our source code run:

```
$ go get github.com/moov-io/iso8583
```

Then, add your (or another user's) fork.

```
$ cd $GOPATH/src/github.com/moov-io/iso8583

$ git remote add $user git@github.com:$user/iso8583.git

$ git fetch $user
```

Now, feel free to branch and push (`git push $user $branch`) to your remote and send us Pull Requests!

## Test Changes

Please run `make check` before submitting your changes. This command runs Go tests and applies various linters.

Since we manage linters centrally across all projects, new checks may be introduced at any time. As a result, your PR might fail due to issues unrelated to your specific changes. We appreciate your help in addressing any issues flagged by new linter checks - either as part of your PR or in a separate one.

## Pull Requests

A good quality PR will have the following characteristics:

* It will be a complete piece of work that adds value in some way.
* It will have a title that reflects the work within, and a summary that helps to understand the context of the change.
* There will be well written commit messages, with well crafted commits that tell the story of the development of this work.
* Ideally it will be small and easy to understand. Single commit PRs are usually easy to submit, review, and merge.
* The code contained within will meet the best practices set by the team wherever possible.
* The code is able to be merged (meaning all CI checks pass)
* A PR does not end at submission though. A code change is not made until it is merged and used in production.

A good PR should be able to flow through a peer review system easily and quickly.
