package main

import (
	brokerManager "github.com/JosephCottingham/mqtt_interface_cli/brokerManager"
	mqtt "github.com/JosephCottingham/mqtt_interface_cli/mqtt"
	ishell "github.com/abiosoft/ishell/v2"
	"os"
	"strconv"
)

func StartShell(brokerData brokerManager.BrokerData, password string) {
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	shell.Set("brokerData", brokerData)
	shell.Set("password", password)

	// display welcome info.
	shell.Println("MQTT Interactive Shell")
	shell.Println("Commands:")
	shell.Println(" -> create\t\t:\tCreate a new broker")
	shell.Println(" -> remove\t\t:\tRemove a given broker")
	shell.Println(" -> ls\t\t\t:\tList")
	shell.Println(" -> connect\t\t:\tConnect to a created broker")
	shell.Println(" -> setpassword\t\t:\tSet a new password")

	shell.AddCmd(&ishell.Cmd{
		Name: "create",
		Help: "Add broker credentials",
		Func: createBroker,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "remove",
		Help: "Remove broker credentials",
		Func: removeBroker,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "ls",
		Help: "List all stored brokers",
		Func: lsBrokers,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "connect",
		Help: "Connect to broker",
		Func: connectListen,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "setpassword",
		Help: "Set password",
		Func: resetPassword,
	})

	// run shell
	shell.Run()
}

func selectBroker(c *ishell.Context) brokerManager.Broker {
	// brokerData := c.Get("brokerData").(brokerManager.BrokerData)
	brokerData, _ := brokerManager.ReadBrokerData(c.Get("password").(string))
	choices := []string{}
	for _, b := range brokerData.Brokers {
		choices = append(choices, b.Name)
	}
	choice := c.MultiChoice(choices, "Select Broker")

	return brokerData.Brokers[choice]
}

func connectListen(c *ishell.Context) {
	broker := selectBroker(c)
	c.Print("Topic: ")
	topic := c.ReadLine()
	mqtt.Connect(broker, topic)
}

func lsBrokers(c *ishell.Context) {
	// brokerData := c.Get("brokerData").(brokerManager.BrokerData)
	brokerData, _ := brokerManager.ReadBrokerData(c.Get("password").(string))
	for i, b := range brokerData.Brokers {
		c.Println(i, ": ", b.Name)
	}
}

func createBroker(c *ishell.Context) {
	broker := brokerManager.Broker{}

	// disable the '>>>' for cleaner same line input.
	c.ShowPrompt(false)
	defer c.ShowPrompt(true) // yes, revert after login.

	// get name of new broker
	c.Print("name: ")
	broker.Name = c.ReadLine()

	// get clientName of new broker
	c.Print("clientName: ")
	broker.ClientName = c.ReadLine()

	// get clientName of new broker
	c.Print("uri: ")
	broker.Uri = c.ReadLine()

	// get clientName of new broker
	c.Print("port: ")
	broker.Port, _ = strconv.Atoi(c.ReadLine())

	// get clientName of new broker
	c.Print("username: ")
	broker.Username = c.ReadLine()

	// get clientName of new broker
	c.Print("password: ")
	broker.Password = c.ReadLine()

	// brokerData := c.Get("brokerData").(brokerManager.BrokerData)
	password := c.Get("password").(string)
	brokerData, _ := brokerManager.ReadBrokerData(password)
	brokerData = brokerManager.AddBroker(brokerData, broker, password)
	c.Set("brokerData", brokerData)

	c.Printf("Broker %s Created\b", broker.Name)
}

func removeBroker(c *ishell.Context) {

	// disable the '>>>' for cleaner same line input.
	c.ShowPrompt(false)
	defer c.ShowPrompt(true) // yes, revert after login.

	broker := selectBroker(c)

	// brokerData := c.Get("brokerData").(brokerManager.BrokerData)
	password := c.Get("password").(string)
	brokerData, _ := brokerManager.ReadBrokerData(password)
	brokerData = brokerManager.RemoveBroker(brokerData, broker, password)
	c.Set("brokerData", brokerData)

	c.Printf("Broker %s Removed\b", broker.Name)
}

func resetPassword(c *ishell.Context) {

	// disable the '>>>' for cleaner same line input.
	c.ShowPrompt(false)
	defer c.ShowPrompt(true) // yes, revert after login.

	// brokerData := c.Get("brokerData").(brokerManager.BrokerData)
	password := c.Get("password").(string)
	brokerData, _ := brokerManager.ReadBrokerData(password)
	c.Printf("Enter New Password: ")
	password = c.ReadPassword()
	c.Printf("Enter New Password: ")
	confirmPassword := c.ReadPassword()
	if password != confirmPassword {
		c.Printf("Entered Passwords Did not Match...")
		return
	}
	brokerManager.WriteBrokerData(brokerData, password)

	c.Printf("Password Reset: Program Reset Required\n")
	os.Exit(4)
}
