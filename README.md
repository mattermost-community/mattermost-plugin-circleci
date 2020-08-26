# Mattermost Plugin CircleCI [![CircleCI branch](https://img.shields.io/circleci/project/github/nathanaelhoun/mattermost-plugin-circleci/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-circleci)

A Work-In-Progress [CircleCI](https://circleci.com) plugin to interact with jobs and builds, with slash commands in Mattermost.

To learn more about plugins, see [the Mattermost plugin documentation](https://developers.mattermost.com/extend/plugins/).

**This plugin is under development and is not ready for production**

## Features (tracking list)

-   [x] Get help

-   [x] Connect to CircleCI, see your profile, disconnect

-   [ ] Setup webhook notifications about successful and failed CircleCI builds

-   [ ] Interact with CircleCI jobs

    -   [ ] Trigger jobs with and without parameters
    -   [ ] Abort a job
    -   [ ] Configure/create/delete jobs
    -   [ ] Get logs from a job in a file attachment, not as a message (this is because the logs can be huge, so it’s easier to preview a file attachment)
    -   [ ] Get artifacts
    -   [ ] Get test results

-   [ ] Meet [requirements](https://developers.mattermost.com/extend/plugins/community-plugin-marketplace/#requirements-for-adding-a-community-plugin-to-the-marketplace) to publish to Marketplace

## Installation

_Coming_

## Development

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. Edit your server configuration as follows:

```json
{
    "ServiceSettings": {
        ...
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    }
}
```

and then deploy your plugin:

```
make deploy
```

You may also customize the Unix socket path:

```
export MM_LOCALSOCKETPATH=/var/tmp/alternate_local.socket
make deploy
```

If developing a plugin with a webapp, watch for changes and deploy those automatically:

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make watch
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with credentials:

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```

## License

Apache License.

## Thanks to

-   Another [CircleCI Plugin](https://github.com/chetanyakan/mattermost-plugin-circleci) by @chetanyakan
-   [Mattermost](https://mattermost.com) for their product
