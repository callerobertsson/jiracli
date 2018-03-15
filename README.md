# Jira CLI

_Simple command line inteface for making Jira queries_

JiraCLI has the following commands

* `issue` - display info on an issue
* `epics` - show epics in a project/team
* `tree` - display a recursive tree of all sub-tasks
* `search` - perform a JQL query

## Project Status

The JiraCLI is still under construction.

## Install

* Clone this repo
* Run `go install`

Copy `example-config.toml` to `$HOME/.jiracli/config.toml` and add Jira
user name and password. 

## Usage

Execute

```jiracli -h```

to display usage information.

/Calle Robertsson, calle@robcon.se

