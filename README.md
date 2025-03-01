# glow-reminder
Glow Reminder offers visual reminders using LEDs controlled via a Telegram bot

The bot is used to create reminders with a choice of colour (green, red or blue) and effect (blinking or static light)

## Configuration

To work with Telegram bot you need to create `.env` file based on `.env.example` and set the value of `TOKEN` variable with your bot token. Additional settings should be made in `config/config.yaml`

## How to start?

```bash
make dev-env-up # runs docker containers with the application
```

Additionally you will need to upload the code in arduino/main.c to your Arduino NodeMCU with wifi module to run the web server on the Arduino

## How to stop?

```bash
make dev-env-stop # stop docker containers
```

