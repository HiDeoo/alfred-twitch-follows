<p align="center">
  <h1 align="center">alfred-workflow-tools 🧰</h1>
</p>

<p align="center">
  <a href="https://github.com/HiDeoo/alfred-workflow-tools/actions/workflows/integration.yml"><img alt="Integration Status" src="https://github.com/HiDeoo/alfred-workflow-tools/actions/workflows/integration.yml/badge.svg"></a>
  <a href="https://github.com/HiDeoo/alfred-workflow-tools/blob/master/LICENSE"><img alt="License" src="https://badgen.now.sh/badge/license/MIT/blue"></a>
  <br /><br />
</p>

`alfred-workflow-tools` is a collection of Go tools that should be invoked by [Alfred](https://www.alfredapp.com/) in [a workflow](https://www.alfredapp.com/help/workflows/).

The results are returned using the [Alfred Script Filter JSON Format](https://www.alfredapp.com/help/workflows/inputs/script-filter/json/) so they can automatically be consumed by a workflow.

## Twitch

`cmd/twitch` fetches Follows (or only Live Follows) for a specific [Twitch](https://twitch.tv) User or a list of Live Streams for a specific Game (and optionally in a specific Language).

### Usage

`cmd/twitch` can be executed in a [Run Script Action](https://www.alfredapp.com/help/workflows/actions/run-script/) to return Twitch Follows:

```shell
$ twitch
{
  "items": [
    {
      "title": "Streamer 1",
      "subtitle": "Science & Technology - 8625 viewers - Coding things",
      "arg": "https://www.twitch.tv/streamer1"
    },
    {
      "title": "Streamer 2",
      "subtitle": "Just Chatting - 2811 viewers - Chatting about things",
      "arg": "https://www.twitch.tv/streamer2"
    },
    …
  ]
}
```

You can also use the `--live` option to only return Twitch Follows that are currently live or alternatively the `--game` option to only return Twitch Live Streams for a specific Game ID (and optionally in a specific Language by using the `--gameLang` option with an ISO 639-1 two-letter language code).

### Configuration

`cmd/twitch` consumes various environment variables that should be [provided by Alfred](https://www.alfredapp.com/help/workflows/advanced/variables/#environment) when invoking the script in a workflow:

| Environment variable | Description                                                    |
| -------------------- | -------------------------------------------------------------- |
| `TWITCH_CLIENT_ID`   | A Twitch application client ID.                                |
| `TWITCH_OAUTH_TOKEN` | A Twitch User access token with the `user:read:follows` scope. |

## BetaSeries

`cmd/betaseries` fetches Shows with Unwatched Episode(s) for a specific [BetaSeries](https://www.betaseries.com) User.

### Usage

`cmd/betaseries` can be executed in a [Run Script Action](https://www.alfredapp.com/help/workflows/actions/run-script/) to return BetaSeries Shows:

```shell
$ betaseries
{
  "items": [
    {
      "title": "Show 1",
      "subtitle": "2 episodes remaining (90 total)",
      "arg": "123456"
    },
    {
      "title": "Show 2",
      "subtitle": "3 episodes remaining (5 total)",
      "arg": "456789"
    },
    …
  ]
}
```

You can also use the `--watched` option to mark all Aired Unwatched Episodes of a BetaSeries Shows as watched.

### Configuration

`cmd/betaseries` consumes various environment variables that should be [provided by Alfred](https://www.alfredapp.com/help/workflows/advanced/variables/#environment) when invoking the script in a workflow:

| Environment variable     | Description                         |
| ------------------------ | ----------------------------------- |
| `BETASERIES_CLIENT_ID`   | A BetaSeries application client ID. |
| `BETASERIES_OAUTH_TOKEN` | A BetaSeries User access token.     |

## GitHub

`cmd/github` fetches Repositories and Recent Contributions for a specific [GitHub](https://github.com/) User. Repositories are cached if the `alfred_workflow_cache` environment variable (automatically provided by Alfred) is set.

### Usage

`cmd/github` can be executed in a [Run Script Action](https://www.alfredapp.com/help/workflows/actions/run-script/) to return GitHub Repositories:

```shell
$ github
{
  "items": [
    {
      "title": "user/repo1",
      "subtitle": "Last activity 2 hours ago",
      "arg": "https://www.github.com/user/repo1"
    },
    {
      "title": "org/repo2",
      "subtitle": "Last activity 5 days ago",
      "arg": "https://www.github.com/org/repo2"
    },
    …
  ]
}
```

You can use the `--clear` option to clear the cached GitHub Repositories.

### Configuration

`cmd/github` consumes an environment variables that should be [provided by Alfred](https://www.alfredapp.com/help/workflows/advanced/variables/#environment) when invoking the script in a workflow:

| Environment variable | Description                     |
| -------------------- | ------------------------------- |
| `GITHUB_OAUTH_TOKEN` | A GitHub personal access token. |

## Contribute

1. [Fork](https://help.github.com/articles/fork-a-repo) & [clone](https://help.github.com/articles/cloning-a-repository) this repository.
1. Build & run the development version using `go run` for the desired tool.

## License

Licensed under the MIT License, Copyright © HiDeoo.

See [LICENSE](https://github.com/HiDeoo/alfred-workflow-tools/blob/master/LICENSE) for more information.
