# Telegram Push Notification Setup

This guide explains how to set up Telegram push notifications for PMail.

## 1. Create a Telegram Bot and Obtain Token

1.  Open Telegram and search for `@BotFather`.
2.  Start a chat with `@BotFather` and send the `/newbot` command.
3.  Follow the instructions to choose a name and a username for your bot.
4.  Upon successful creation, `@BotFather` will provide you with an HTTP API token. This is your `tgBotToken`.

## 2. Obtain your Telegram Chat ID

To send messages to your Telegram account, the bot needs your `chatId`.

1.  Start a conversation with your newly created bot by sending any message.
2.  Open your web browser and navigate to `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates` (replace `<YOUR_BOT_TOKEN>` with your actual bot token).
3.  Look for the `"chat"` object in the JSON response. Your `chatId` will be the value associated with the `"id"` key within that object.

    Example JSON snippet:
    ```json
    {
      "update_id": 123456789,
      "message": {
        "message_id": 123,
        "from": { ... },
        "chat": {
          "id": 1234567890, // This is your chatId
          "first_name": "YourName",
          "type": "private"
        },
        "date": 1678901234,
        "text": "Hello Bot"
      }
    }
    ```

## 3. Configure PMail for Telegram Push

1.  Copy the Telegram push plugin binary to the `/plugins` directory of your PMail installation.

2.  Add the following configuration to your `config.json` file (located in the `config` directory relative to your PMail executable):

    ```jsonc
    {
      "tgChatId": "YOUR_CHAT_ID",     // Replace with your Telegram Chat ID
      "tgBotToken": "YOUR_BOT_TOKEN"  // Replace with your Telegram Bot Token
    }
    ```

3.  Ensure that the Telegram push plugin is enabled in your PMail configuration. (Further details on enabling plugins will be provided in the main configuration documentation if not already present).

After configuring, restart your PMail server for the changes to take effect.
