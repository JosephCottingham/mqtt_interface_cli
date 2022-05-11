# MQTT CLI Interactive Shell

This is a CLI tool for interfacing with MQTT brokers. Unlike other MQTT CLI tools, this software operates in a shell-based environment with stored credentials (encrypted) designed to provide quick access to the brokers you as a developer frequently use.


## Shell Commands:
    - create
        - Used to create a new broker for which you may connect to.
    - remove
        - Remove undesired stored brokers.
    - ls
        - List all stored brokers
    - connect
        - Connect to a given broker and subscribe to a topic
    - setpassword
        - Set the user password for which your credentials are encrypted with (Recommended that you use a length divisible by 16 so no padding is used.)



## How To Run:

```
go run main.go
```