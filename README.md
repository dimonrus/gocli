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

1. Define stage config(global.yaml and local.yaml with depends on global.yaml)
    _global.yaml_
    ``` 
    depends:
    project:
      name: dna
      debug: false
    web:
      port: 8080
      host: 0.0.0.0
    arguments:
      app:
        type: string
        label: application type
      script:
        type: string
        label: console command name
      consumer:
        type: string
        label: consumer name
      count:
        type: int
        label: count
      success:
        type: bool
        lable: success
      part:
        type: float
        label: percent
    ```
   _local.yaml_
    ```
    depends: global
      project:
        debug: true
    ```
2. Define golang config file with structure for defined config
    ``` 
    type Config struct {
        Project struct {
            Name  string
            Debug bool
        }
        Web struct {
            Port int
            Host string
        }
        Arguments gocli.Arguments
    }
    ```
3. Init application
    ```
    var config Config
    environment := os.Getenv("ENV")
    if environment == "" {
       environment = "local"
    }
    rootPath, err := filepath.Abs("")
    if err != nil {
        panic(err)
    }
    app := gocli.NewApplication(environment, rootPath+"/config/yaml", &config)
    app.ParseFlags(&config.Arguments)

    appType, ok := config.Arguments["app"]
    if ok != true {
        app.FatalError(errors.New("app type is not presents"))
    }
    ```
4. Listen command port
    ```
    exit := make(chan, struct{})
    go func() {
        err = app.Start(":3333", func(command *gocli.Command) {
            v := command.Arguments()[0]
            app.SuccessMessage("Receive command: "+command.String(), command)
            if v.Name == "exit" {
                app.AttentionMessage("Exit...", command)
                exit <- struct{}{}
            } else if v.Name == "show" {
                app.AttentionMessage(gohelp.AnsiYellow+"The show is began"+gohelp.AnsiReset, command)
            } else {
                app.AttentionMessage(gohelp.AnsiRed+"Unknown command: "+command.String()+gohelp.AnsiReset, command)
            }
        })
    }()
    <- exit
    ```

#### If you find this project useful or want to support the author, you can send tokens to any of these wallets
- Bitcoin: bc1qgx5c3n7q26qv0tngculjz0g78u6mzavy2vg3tf
- Ethereum: 0x62812cb089E0df31347ca32A1610019537bbFe0D
- Dogecoin: DET7fbNzZftp4sGRrBehfVRoi97RiPKajV