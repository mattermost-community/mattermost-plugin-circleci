## Configuration

### Step 1: Configure the bot account in Mattermost

If you have an existing Mattermost user account with the name circleci, the plugin will post using the circleci account but without a BOT tag.

To prevent this, either:

- Convert the circleci user to a bot account by running `mattermost user convert circleci --bot` in the Mattermost CLI.

or

- If the user is an existing user account you want to preserve, change its username and restart the Mattermost server. Once restarted, the plugin will create a bot account with the name `circleci`.

### Step 2: Configure the plugin in Mattermost

To generate the keys needed, go to **System Console > Plugins > CircleCI**:

1. Generate a new value for **Webhooks Secret**. If the generated secret contains a forwardslash, please regenerate it.
2. Generate a new value for **At Rest Encryption Key**.
3. Select **Save**.
4. Go to **System Console > Plugins > Management** and choose **Enable** to enable the CircleCI plugin.

You're all set!
