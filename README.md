# Twitter Tools

[![Matrix](https://img.shields.io/matrix/054h509j4h509hjtrj455g:matrix.org?label=chat&logo=Matrix)](https://matrix.to/#/#054h509j4h509hjtrj455g:matrix.org)


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

## Flow

Flow actions are defined in your `config.yml`. If you don't already have a local config file then an example one will be created for you on first execution. You can always check the latest _example config_ [here](https://github.com/tupini07/twitter-tools/blob/main/app_config/config.example.yml). The example config contains all the possible steps that can be used as part of the flow, which are described in more depth below.

### Common config

```yaml
flow:
  repeat: true # repeat flow execution as soon as it completes
  max_total_following: 4500 # no users will be followed if you're already following this number or more
  steps:
    - list of steps we want to execute
```

### Possible Steps

TODO all below

#### Follow followers of others

```yaml
- follow_followers_of_others:
    max_to_follow: (number) the maximum number of people we want to follow as part of this step. If not provided then there will be no follow limit
    max_sources_to_pick: (number) pick this number of sources randomly from the list of "others". If not provided then the whole list will be used
    others: # screen names of other accounts to follow followers from
      - list of twitter handlers
```

#### Follow all your followers

#### Unfollow people that are not following you

#### Wait

#### Wait until day

#### Random step


## TODOs

- Save database and config to a default location
  - Maybe `~/.local/share/twitter-tools/` ?
