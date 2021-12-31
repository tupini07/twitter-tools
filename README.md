# Twitter Tools

This project is a small collection of tools that automate certain tasks for Twitter. I used it to teach myself the [Go programming language](https://github.com/golang/go) along the way.

```
NAME:                                                                                        
   twitter-tools - Collection of tools to manage a Twitter account                           
                                                                                             
USAGE:                                                                                       
   twitter-tools.exe [global options] command [command options] [arguments...]               
                                                                                             
VERSION:                                                                                     
   0.0.3                                                                                     
                                                                                             
COMMANDS:                                                                                    
   unfollow-bad-friends, ubf        Unfollows bad friends starting from the oldest friendship
   follow-all-followers, faf        Ensure followers of the current user are being followed  
   follow-followers-of-other, ffoo  Follow followers of other(s) users                       
   do-flow, df                      Performs the flow actions defined in config.yml          
   help, h                          Shows a list of commands or help for one command         
                                                                                             
GLOBAL OPTIONS:                                                                              
   --help, -h     show help (default: false)                                                 
   --version, -v  print the version (default: false)                                         
```

## Supported tasks

TODO

## Flow

TODO

## TODOs

- Save database and config to a default location
  - Maybe `~/.local/share/twitter-tools/` ?