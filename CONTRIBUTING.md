# Contributing to `Topfew`

Topfew is hosted in this GitHub repository 
at `github.com/timbray/topfew` and welcomes 
contributions.

As of mid-2024, Topfew is probably more or
less complete. It is well-tested. Its performance at
processing streams can keep up with most streams
and it is dramatically faster when processing files,
where it processes multiple segments in parallel.

If you disagree and want to contribute to Topfew,
the first step in making a change is typically to 
raise an Issue to allow for discussion of the idea. 
This is important because possibly Topfew already
does what you want, in which case perhaps what’s 
needed is a documentation fix. Possibly the idea 
has been raised before but failed to convince Topfew’s
maintainers. (Doesn’t mean it won’t find favor now;
times change.)

Assuming there is agreement that a change in Topfew
is a good idea, the mechanics of forking the repository,
committing changes, and submitting a pull request are
well-described in many places; there is nothing 
unusual about Topfew.

### Code Style

The coding style suggested by the Go community is 
used in Topfew. See the
[style doc](https://go.dev/wiki/CodeReviewComments) for details.

Try to limit column width to 120 characters for both code and markdown documents
such as this one.

### Signing commits

Commits should be signed (not just the `-s` “signd off on”) with
any of the [styles GitHub supports](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits).
Note that you can use `git config` to arrange that your commits are
automatically signed with the right key.

### Running Tests

In the repository root `go test ./...` runs unit tests
with all the defaults, which is a decent check for basic
sanity and correctness.

The following command runs the Go linter; submissions 
need to be free of lint errors.

```shell
golangci-lint run  
```
