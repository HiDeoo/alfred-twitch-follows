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

`cmd/twitch` fetches Follows (or only Live Follows) for a specific [Twitch](https://twitch.tv) User.

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

You can also use the `--live` option to only return Twitch Follows that are currently live.

### Configuration

`cmd/twitch` consumes various environment variables that should be [provided by Alfred](https://www.alfredapp.com/help/workflows/advanced/variables/#environment) when invoking the script in a workflow:

| Environment variable | Description                                                    |
| -------------------- | -------------------------------------------------------------- |
| `TWITCH_CLIENT_ID`   | A Twitch application client ID.                                |
| `TWITCH_OAUTH_TOKEN` | A Twitch User access token with the `user:read:follows` scope. |

## Contribute

1. [Fork](https://help.github.com/articles/fork-a-repo) & [clone](https://help.github.com/articles/cloning-a-repository) this repository.
1. Build & run the development version using `go run` for the desired tool.

## License

Licensed under the MIT License, Copyright © HiDeoo.

See [LICENSE](https://github.com/HiDeoo/alfred-workflow-tools/blob/master/LICENSE) for more information.
