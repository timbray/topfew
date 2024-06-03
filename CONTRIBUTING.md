# Contributing to Topfew

Topfew is hosted in this GitHub repository 
at `github.com/timbray/topfew` and welcomes 
contributions.

This is release 2.0 of Topfew, which is probably more 
or less complete. It is well-tested. Its performance
at processing streams can keep up with most streams
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

### Format of the Commit Message

We follow the conventions described in [How to Write a Git Commit
Message](http://chris.beams.io/posts/git-commit/).

Be sure to include any related GitHub issue references in the commit message,
e.g. `Closes: #<number>`.

The [`CHANGELOG.md`](./CHANGELOG.md) and release page uses **commit message
prefixes** for grouping and highlighting. A commit message that
starts with `[prefix:] ` will place this commit under the respective
section in the `CHANGELOG`.
- `chore:` - Use for repository related activities
- `fix:` - Use for bug fixes
- `docs:` - Use for changes to the documentation
- `kaizen:` - Use for improvements, including optimization and new features

If your contribution falls into multiple categories, e.g. `api` and `pat` it
is recommended to break up your commits using distinct prefixes.

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
