# Contributing

Thank you for thinking of contributing!

## You want to report a bug

Thank you for taking the time to report a bug!  I'm sure others are dealing with the bug, so you are helping others by
reporting.

Please [file an issue](https://github.com/halprin/delete-dynamodb-items/issues/new/choose).
Provide...

- A meaningful title.
- What you expected.
- What actually happened.

Bonus points if you can provide a script of `aws` CLI commands and `delete-dynamodb-items` that results in the bug.

## You want to suggest a new feature

I appreciate you wanting to improve the functionality!  I bet others would appreciate the new feature too.

Please [file an issue](https://github.com/halprin/delete-dynamodb-items/issues/new/choose).
Provide...

- A meaningful title.
- A description of the feature.

## You want to fix a bug or create a new feature

That's great!  You can see a list of known bugs and features suggestions under
[issues](https://github.com/halprin/delete-dynamodb-items/issues).  When you pick one, leave a comment that you'd like
work on it and any thoughts you have on it.  I'll assign you, and you're off to the races!  Create a fork and then a PR.

If you want to fix an unreported bug, I would appreciate a bug report first so that I can be made aware of the bug as
soon as possible.  Feel free to say in the description or a subsequent comment that you'd like to work on it.

If you want to work on an unsuggested feature, please file a feature suggestion first so we can discuss the feature
before charging ahead.  Feel free to say in the description or a subsequent comment that you'd like to work on it.

### Development Dependencies

- [GoLang](https://go.dev).
- [Bash](https://www.gnu.org/software/bash/).
- [Docker](https://www.docker.com).
- [AWS CLI](https://aws.amazon.com/cli/).

### Building the Application and Running Tests

As part of modifying the code, you'll want to test your changes.

Run the following to compile your own copy from source.

```shell
make compile
```

Run the following to execute all of the tests.

```shell
make test
```

Don't forget to add or modify tests as you modify the source code.  Tests fall into two categories: unit tests or
integration tests.  Unit tests reside next to the code it tests, a GoLang norm.  The integration tests are defined in
the [`Makefile`](./Makefile) for now.
