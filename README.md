# gocli
Entry point for application

Features
1. Support environment configs with dependencies.
    _Each yaml config file may contain depends key to preload configs from parent environment config files_
2. Support argument customization via config
    _You can define your own flags parsed automatically on application starts_
3. Support processing commands through socket connection
    _You can define your own command passed through socket connection. Command implementation_
4. Logger interface
    Minimal function set for basic logger
5. Standard application instance (DNApp) out of the box.

# Usage
