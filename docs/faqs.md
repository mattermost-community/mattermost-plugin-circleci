# Frequently Asked Questions

## How does the plugin save user data for each connected CircleCI user?

CircleCI user tokens are AES encrypted with an At Rest Encryption Key configured in the plugin's settings page. Once encrypted, the tokens are saved in the `PluginKeyValueStore` table in your Mattermost database.
