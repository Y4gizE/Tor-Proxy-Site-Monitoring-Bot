
- **SocksPort**: This sets the Tor proxy to listen on `port 9150` on all network interfaces (`0.0.0.0`).
- **DataDirectory**: Specifies where Tor stores its data. Update the path according to your system.

## Setup Instructions

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/tor-site-monitor-bot.git
    cd tor-site-monitor-bot
    ```

2. Set your Telegram bot token and chat ID in the code by updating the constants:

    ```go
    const telegramAPI = "https://api.telegram.org/bot<your_bot_token>/sendMessage"
    const chatID = "<your_chat_id>"
    ```

3. Run the application:

    ```bash
    go run main.go
    ```

## How It Works

- The bot checks if Tor is working by making a request to `http://check.torproject.org`.
- Every 5 minutes, the bot sends a GET request to the target site through the Tor network.
- If the site is down or the response time exceeds 100 ms, the bot sends a warning message to your Telegram group.
- If the site is up and the response time is acceptable, the bot reports the status along with the response time.

## Customization

- **Target Site**: You can change the site to monitor by modifying the `targetSite` constant:

    ```go
    const targetSite = "https://example.com"
    ```

- **Response Time Threshold**: You can adjust the warning threshold by updating the `slowResponseThreshold` constant:

    ```go
    const slowResponseThreshold = 100 * time.Millisecond
    ```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
