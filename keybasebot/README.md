# Keybase Bot for Home Assistant

This project is a WIP. This is a prototype/proof of concept. It has not been tested extensively and not suitable for production use-cases.

This has only been compiled for and tested on armv7 (Raspberry Pi) but can easily be compiled for other platforms.

Currently this addon allows messages to be sent from Home Assitant to Keybase. To set up this integration, you must geenerate an account for your bot with a paper key. The bot can currently post messages to a single channel for a single team. The bot must be a member of this team.

## Configuration

Key | Description 
----|------------
username | username for the bot
paperKey | a paper key to authenticate the bot
teamName | the bot will send messages to the default channel for this team

## Usage

To send messages using automation, use event `NOTIFY_KEYBASE` and provide a `message` attribute in the data section. For example,
```
action:
- event: NOTIFY_KEYBASE
      event_data:
        message: "This is a test!"
```

## Build Instructions

For testing

1. `make build`
2. Upload Dockerfile, config.json, keybasebot-armv7 to the addons folder on your Home Assistant server (eg, on hass.io: `/addons/castlebot`)
3. Refresh (top right) on the Add-on Store tab (in Hass.io)
4. Install & start the addon
