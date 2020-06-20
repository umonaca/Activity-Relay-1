# Activity Relay Server - `Project Improve`

## Yet another powerful customizable ActivityPub relay server written in Go.

[![GitHub Actions](https://github.com/yukimochi/Activity-Relay/workflows/Test/badge.svg?branch=project_improve)](https://github.com/yukimochi/Activity-Relay/tree/project_improve)
[![codecov](https://codecov.io/gh/yukimochi/Activity-Relay/branch/project_improve/graph/badge.svg)](https://codecov.io/gh/yukimochi/Activity-Relay/branch/project_improve)

![Connect to Fediverse, join Fediverse.](misc/title.png)

## Note: `Project Improve` is Work in Progress. **NOT WORKS**

## Packages

 - `github.com/yukimochi/Activity-Relay`
 - `github.com/yukimochi/Activity-Relay/api`
 - `github.com/yukimochi/Activity-Relay/command`
 - `github.com/yukimochi/Activity-Relay/deliver`
 - `github.com/yukimochi/Activity-Relay/models`

## Requirement

 - [Redis](https://github.com/antirez/redis)

## Run

### API Server

```bash
Activity-Relay -c <Path of config file> server
```

### Job Worker

```bash
Activity-Relay -c <Path of config file> worker
```

### CLI Management Utility

```bash
Activity-Relay -c <Path of config file> control
```

## Config

### YAML Format

```yaml config.yml
ACTOR_PEM: actor.pem
REDIS_URL: redis://localhost:6379

RELAY_BIND: 0.0.0.0:8080
RELAY_DOMAIN: relay.toot.yukimochi.jp
RELAY_SERVICENAME: YUKIMOCHI Toot Relay Service
RELAY_SUMMARY: |
  YUKIMOCHI Toot Relay Service is Running by Activity-Relay
RELAY_ICON: https://example.com/example_icon.png
RELAY_IMAGE: https://example.com/example_image.png
```

### Environment Variable

This is **Optional** : When config file not exist, use environment variables.

 - `ACTOR_PEM`
 - `REDIS_URL`
 - `RELAY_BIND`
 - `RELAY_DOMAIN`
 - `RELAY_SERVICENAME`
 - `RELAY_SUMMARY`
 - `RELAY_ICON`
 - `RELAY_IMAGE`

## Project Sponsors

Thank you for your support.

### Monthly Donation

**[My Doner List](https://relay.toot.yukimochi.jp#patreon-list)**
  
#### Donation Platform
 - [Patreon](https://www.patreon.com/yukimochi)
 - [pixiv fanbox](https://yukimochi.fanbox.cc)
 - [fantia](https://fantia.jp/fanclubs/11264)

### IDE Support

[![Jetbrains Logo](misc/jetbrains.svg)](https://www.jetbrains.com/?from=Activity-Relay)

[JetBrains Free License Programs for Open Source](https://www.jetbrains.com/?from=Activity-Relay)

