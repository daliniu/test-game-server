# micro bot

The **micro bot** is a bot that sits inside your microservices environment which you can interact with via Slack, HipChat, XMPP, etc. 
It mimics the functions of the CLI via messaging.

<p align="center">
  <img src="https://github.com/micro/docs/blob/master/images/bot.png" />
</p>

## Supported Inputs

- Slack
- HipChat

## Getting Started

### Install Micro

```go
go get github.com/micro/micro
```

### Run with Slack

```shell
micro bot --inputs=slack --slack_token=SLACK_TOKEN
```

<img src="https://github.com/micro/docs/blob/master/images/slack.png">
-

### Run with HipChat

```shell
micro bot --inputs=hipchat --hipchat_username=XMPP_USER --hipchat_password=XMPP_PASSWORD
```

<img src="https://github.com/micro/docs/blob/master/images/hipchat.png">
-

Use multiple inputs by specifying a comma separated list

```shell
micro bot --inputs=hipchat,slack --slack_token=SLACK_TOKEN --hipchat_username=XMPP_USER --hipchat_password=XMPP_PASSWORD
```

### Help

In slack
```shell
micro help

deregister service [definition] - Deregisters a service
echo [text] - Returns the [text]
get service [name] - Returns a registered service
health [service] - Returns health of a service
hello - Returns a greeting
list services - Returns a list of registered services
ping - Returns pong
query [service] [method] [request] - Returns the response for a service query
register service [definition] - Registers a service
the three laws - Returns the three laws of robotics
time - Returns the server time
```

## Adding new Commands

Commands are functions executed by the bot based on text based pattern matching.

### Write a Command

```go
import "github.com/micro/micro/bot/command"

func Ping() command.Command {
	usage := "ping"
	description := "Returns pong"

	return command.NewCommand("ping", usage, desc, func(args ...string) ([]byte, error) {
		return []byte("pong"), nil
	})
}
```

### Register the command

Add the command to the Commands map with a pattern key that can be matched by golang/regexp.Match

```go
import "github.com/micro/micro/bot/command"

func init() {
	command.Commands["^ping$"] = Ping()
}
```

### Link the Command

Drop a link to your command into the top level dir

link_command.go:
```go
import _ "path/to/import"
```

## Adding new Inputs

Inputs are plugins for communication e.g Slack, HipChat, XMPP, IRC, SMTP, etc, etc. 

New inputs can be added in the following way.

### Write an Input

Write an input that satisfies the Input interface.

```go
type Input interface {
	// Provide cli flags
	Flags() []cli.Flag
	// Initialise input using cli context
	Init(*cli.Context) error
	// Stream events from the input
	Stream() (Conn, error)
	// Start the input
	Start() error
	// Stop the input
	Stop() error
	// name of the input
	String() string
}
```

### Register the input

Add the input to the Inputs map.

```go
import "github.com/micro/micro/bot/input"

func init() {
	input.Inputs["name"] = MyInput
}
```

### Link the input

Drop a link to your input into the top level dir

link_input.go:
```go
import _ "path/to/import"
```

## Commands as Services

The micro bot supports the ability to create commands as microservices. 

### How does it work?

The bot watches the service registry for services with it's namespace. The default namespace is `go.micro.bot`. 
Any service within this namespace will automatically be added to the list of available commands. When a command 
is executed, the bot will call the service with method `Command.Exec`. It also expects the method `Command.Help` 
to exist for usage info.


The service interface is as follows and can be found at [micro/bot/proto](https://github.com/micro/micro/blob/master/bot/proto/bot.proto)

```
syntax = "proto3";

package go.micro.bot;

service Command {
	rpc Help(HelpRequest) returns (HelpResponse) {};
	rpc Exec(ExecRequest) returns (ExecResponse) {};
}

message HelpRequest {
}

message HelpResponse {
	string usage = 1;
	string description = 2;
}

message ExecRequest {
	repeated string args = 1;
}

message ExecResponse {
	bytes result = 1;
	string error = 2;
}
```

### Example

Here's an example echo command as a microservice

```go
package main

import (
	"fmt"
	"strings"

	"github.com/micro/go-micro"
	"golang.org/x/net/context"

	proto "github.com/micro/micro/bot/proto"
)

type Command struct{}

// Help returns the command usage
func (c *Command) Help(ctx context.Context, req *proto.HelpRequest, rsp *proto.HelpResponse) error {
	// Usage should include the name of the command
	rsp.Usage = "echo"
	rsp.Description = "This is an example bot command as a micro service which echos the message"
	return nil
}

// Exec executes the command
func (c *Command) Exec(ctx context.Context, req *proto.ExecRequest, rsp *proto.ExecResponse) error {
	rsp.Result = []byte(strings.Join(req.Args, " "))
	// rsp.Error could be set to return an error instead
	// the function error would only be used for service level issues
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.bot.echo"),
	)

	service.Init()

	proto.RegisterCommandHandler(service.Server(), new(Command))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
```