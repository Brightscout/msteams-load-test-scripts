# MS Teams load test scripts

MS Teams load-test-scripts provides a set of scripts written in [Go](https://golang.org/) and [JS](https://developer.mozilla.org/en-US/docs/Web/JavaScript) to help profiling MS Teams under heavy load, simulating real-world usage of a server installation at scale.

## Setup

Make sure you have the following components installed:  

- Go - v1.18 - [Getting Started](https://golang.org/doc/install)
    > **Note:** If you have installed 'Go' to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).

- Install the K6 load testing tool from [here](https://k6.io/docs/get-started/installation).

- Clone the repo using the command:
    ```
    git clone git@github.com:Brightscout/msteams-load-test-scripts.git
    ``` 
    or 
    ```
    git clone https://github.com/Brightscout/msteams-load-test-scripts.git
    ```

## How to use
- Create a `config.json` file.
    - Run command to copy the sample `config.sample.json` file.
    ```
    cp config/config.sample.json config/config.json
    ```
    - Configure the `config.json` file created according to the load to be tested.
    - Go to config [docs](docs/config.md) to check details on the config settings.

- Run the command `make build` to create a new binary file for the load test scripts.

- Run the command `make init_users` to initialize users present in the config file.

- Run the command `make create_channels` to create the new MS Teams channels and add the above users to them. 

- Go to your Mattermost account and connect any Mattermost channel with these new MS Teams channels.

- Run the command `make create_chats` to create MS Teams chats between the users.

- Run the command `make create_posts` to create the random posts in the MS Teams channels and chats.

- Run the command `make clear_store` to clear all the stored data present in the temporary file called `temp_store.json` to start load testing with new details.
