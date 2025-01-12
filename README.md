# Global Entry Appointment Bot

Global Entry Appointment Bot is a basic Command Line Interface (CLI) tool designed to help you monitor and get a Global Entry Appointment. 

Sure there are services that do this for you at a cost but why do that when you host your own homelab?

## Features
* Can monitor multiple Global Entry Appointment Locations
* Can send discord message to notify when an appointment is available
* Can run on a cron schedule within the binary (no need for external cron)

## Getting Started

### Installation

#### Binary
```shell
# TODO
```

#### Docker
```shell
# TODO
```

### Usage


#### Lookup
We can Look Up airport IDs by using the city and/or state:
```shell
# Help
TZ="America/New_York" ./global-entry-appointment-bot lookup --help
Lookup location Information for global-entry-appointment-bot

Usage:
  global-entry-appointment-bot lookup [flags]

Flags:
  -c, --city string    The city to sort lookups on
  -h, --help           help for lookup
  -s, --state string   The state to sort lookups on

Global Flags:
      --config string       Path to the configuration file
  -f, --log-format string   Log format (text, json) (default "text")
  -l, --log-level string    Log level (trace, debug, info, warn, error, fatal, panic) (default "info")
```
```shell
# Example
TZ="America/New_York" ./global-entry-appointment-bot lookup  --state "AK"
INFO[2025-01-12T10:48:42-05:00] Appointment Location Found                    Address="4600 Postmark Drive, RM NA 207" City="Anchorage " ID=7540 Name="Anchorage Enrollment Center" State=AK
INFO[2025-01-12T10:48:42-05:00] Appointment Location Found                    Address="Room 1320A" City=Fairbanks ID=14381 Name="Fairbanks Enrollment Center" State=AK
```

#### Run
We can run it one-off or on a cron schedule with and without Discord notifications:
```shell
# Help
TZ="America/New_York" ./global-entry-appointment-bot run --help
Run the global-entry-appointment-bot

Usage:
  global-entry-appointment-bot run [flags]

Flags:
  -c, --cron-expression string     The cron schedule to run the command
      --discord-token string       Discord Notification Token
      --discord-user-id string     Discord User ID
  -h, --help                       help for run
  -n, --limit string               The number of results returned by the query (useful when trying to narrow down multiple results) (default "1")
  -i, --location-ids stringArray   The Location IDs where you'd like your appointment (hint: use the lookup command to find your ID)
  -m, --minimum string             The number of minimum available appointments (default "1")
  -b, --order-by string            How to order the results (default "soonest")

Global Flags:
      --config string       Path to the configuration file
  -f, --log-format string   Log format (text, json) (default "text")
  -l, --log-level string    Log level (trace, debug, info, warn, error, fatal, panic) (default "info")
```
```shell
# Run Once
TZ="America/New_York" ./global-entry-appointment-bot run --location-ids 7540
INFO[2025-01-12T10:49:27-05:00] Appointment Found                             Appointment="&{7540 2025-01-14 03:50:00 -0500 EST 2025-01-14 04:00:00 -0500 EST true 10 false}"
INFO[2025-01-12T10:49:28-05:00] Appointment Information                       message="# Appointment Detected!\n* **Date:** 2025-01-14\n* **Time:** 03:50 AM\n* **Name:** Anchorage Enrollment Center\n* **State:** AK\n* **City:** Anchorage \n* **Address:** 4600 Postmark Drive, RM NA 207"
```

```shell
# Run on a Cron Schedule
TZ="America/New_York" ./global-entry-appointment-bot run --location-ids 7540 --cron-expression "@every 10s"
INFO[2025-01-12T10:50:41-05:00] Cron mode                                     cron="&{@every 10s}"
INFO[2025-01-12T10:50:51-05:00] Appointment Found                             Appointment="&{7540 2025-01-14 03:50:00 -0500 EST 2025-01-14 04:00:00 -0500 EST true 10 false}"
INFO[2025-01-12T10:50:51-05:00] Appointment Information                       message="# Appointment Detected!\n* **Date:** 2025-01-14\n* **Time:** 03:50 AM\n* **Name:** Anchorage Enrollment Center\n* **State:** AK\n* **City:** Anchorage \n* **Address:** 4600 Postmark Drive, RM NA 207"
INFO[2025-01-12T10:51:01-05:00] Appointment Found                             Appointment="&{7540 2025-01-14 03:50:00 -0500 EST 2025-01-14 04:00:00 -0500 EST true 10 false}"
INFO[2025-01-12T10:51:01-05:00] Appointment Information                       message="# Appointment Detected!\n* **Date:** 2025-01-14\n* **Time:** 03:50 AM\n* **Name:** Anchorage Enrollment Center\n* **State:** AK\n* **City:** Anchorage \n* **Address:** 4600 Postmark Drive, RM NA 207"
```

#### Configuration Variables, Files and flags.
| Env Variable            | Config File Variable    | Description                                                                                                                                                                            |
|-------------------------|-------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `TZ`                    | `N/A`                   | Standard `TZ` variable from [Wikipedia](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones). This ensures that the appointment time is correctly converted to your timezone. |
| `LOGGING_LEVEL`         | `logging.level`         | Log level (trace, debug, info, warn, error, fatal, panic) (default "info")                                                                                                             |
| `LOGGING_FORMAT`        | `logging.format`        | Log format (text, json) (default "text")                                                                                                                                               |
| `CRON_EXPRESSION`       | `cron.Expression`       | The cron schedule to run the command                                                                                                                                                   |
| `RUN_LIMIT`             | `run.Limit`             | The number of results returned by the query (useful when trying to narrow down multiple results) (default "1")                                                                         |
| `N/A`                   | `run.LocationIds`       | The Location IDs where you'd like your appointment (hint: use the lookup command to find your ID)                                                                                      |
| `RUN_MINIMUM`           | `run.Minimum`           | The number of minimum available appointments (default "1")                                                                                                                             |
| `RUN_ORDERBY`           | `run.OrderBy`           | How to order the results (default "soonest")                                                                                                                                           |
| `NOTIFY_DISCORD_TOKEN`  | `notify.discord.Token`  | Discord Notification Token                                                                                                                                                             |
| `NOTIFY_DISCORD_USERID` | `notify.discord.UserID` | Discord User ID                                                                                                                                                                        |

##### Example Configuration File

Configuration file should support any format supported by [spf13/viper](https://github.com/spf13/viper).

```yaml
# config.yaml
---
logging:
  level: "trace"
  format: "text"
cron:
  Expression: "@every 5m"
run:
  Limit: 10
  LocationIds:
    - "7540"
  Minimum: 2
  OrderBy: "soonest"
notify:
  discord:
    Token: ""
    UserID: ""
```
