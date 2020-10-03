# Contributing to this project

### I saw a bug, I have a feature request or a suggestion

Please fill a [Github Issue](https://github.com/nathanaelhoun/mattermost-plugin-circleci/issues/new/choose), it will be very useful!

### I want to code

Pull Requests are welcome! You can contact me on the [Mattermost Community ~plugin-circleci channel](https://community.mattermost.com/core/channels/plugin-circleci) where I am `@nathanaelhoun`.

## Development

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. After configuring it, just run:

```
make deploy
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```
