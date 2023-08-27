# DiscordBot

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

## About The Project

There are many great Discord bots out there, but I wanted to develop my own. I wanted the freedom of being able to easily add my own commands and tailor it to my own needs. After thinking a bit about which language to use, I decided on developing it in Go. I decided on Go after my newfound love for the programming language after the cloud technologies course I had during my bachelor in programming. Of course this Discord bot may not cover all of your own needs, but you're more than welcome to fork this repo and change it how you like.

### Built With

This discord bot is built with:

* [Go version 1.20](https://golang.org/)
* [DiscordGo](https://github.com/bwmarrin/discordgo)
* [dca](https://github.com/jonas747/dca)
* [Docker and docker compose v2](https://www.docker.com)

## Getting Started

To get a local copy up and running follow these simple example steps.

### Prerequisites

In order to run this project you will need the following:

* Go 1.20 installed
* Docker and Docker Compose v2 installed
* Discord account
* Google developer account

### Installation

1. Go to the [Discord developer portal](https://discord.com/developers)
2. Create a new application
3. Add a bot user to the application
4. Get the token for the bot
5. Clone the repository

   ```sh
   git clone https://github.com/Dahl99/DiscordBot.git
   ```

6. Install dependencies

   ```sh
   go mod download all
   ```

7. Create the environment variables file `.env` from the `.env.example` file in the root folder and add the following:

    ```dotenv
    APP_ENVIRONMENT=LOCAL or TEST or PRODUCTION
    BOT_PREFIX="Your prefix for bot commands"
    BOT_STATUS="Your bot status message"
    BOT_GUILD_JOIN_MESSAGE="Your bot guild join message"
    DISCORD_TOKEN=Your discord token
    YOUTUBE_API_KEY=Your youtube api key
    ```

## Usage

To run the discord bot from root directory, execute the following command
```shell
go run main.go
```

Or to start the docker containers using docker compose V2

```sh
docker compose up
```

## Contact

Project Link: [https://github.com/Dahl99/DiscordBot](https://github.com/Dahl99/DiscordBot)
