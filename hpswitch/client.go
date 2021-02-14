package hpswitch

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

type client struct {
	Host     string
	User     string
	Password string
}

func (c client) readVlan(vlan string) (string, string, []map[string]string) {

	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", c.Host+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run("sys\ndisplay vlan " + vlan + "\nquit"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}

	sshOut := b.String()

	if strings.Contains(sshOut, "This VLAN does not exist.") {
		return "", "", nil
	}

	out := strings.Split(sshOut, "display vlan "+vlan)[1]
	out = strings.Split(out, "[")[0]

	vlanID := strings.Split(out, "VLAN ID: ")[1]
	vlanID = strings.Split(vlanID, "\n")[0]
	vlanID = strings.Replace(vlanID, "\r", "", -1)

	description := strings.Split(out, "Description: ")[1]
	description = strings.Split(description, "\n")[0]
	description = strings.Replace(description, "\r", "", -1)

	taggedPorts := strings.Split(out, "Tagged ports:")[1]
	taggedPorts = strings.Split(taggedPorts, "Untagged ports:")[0]
	taggedPortsList := strings.Fields(taggedPorts)

	taggedPortsMaps := make([]map[string]string, 0)

	for i := 0; i < len(taggedPortsList); i++ {
		if taggedPortsList[i] != "" {
			m := make(map[string]string)
			m["port"] = taggedPortsList[i]
			taggedPortsMaps = append(taggedPortsMaps, m)
		}
	}

	return vlanID, description, taggedPortsMaps
}

func (c client) setVlan(vlan string, description string, taggedPorts []interface{}) string {

	commandString := "sys\nvlan " + vlan + "\n" + "description " + description + "\n"

	fmt.Println(commandString)

	for _, taggedPortIF := range taggedPorts {
		taggedPortMap := taggedPortIF.(map[string]interface{})
		taggedPort := taggedPortMap["port"].(string)
		commandString = commandString + "interface " + taggedPort + "\nport trunk permit vlan " + vlan + " " + vlan + "\n"
	}

	commandString = commandString + "quit\nsave f"

	// Establish connection to switch
	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", c.Host+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(commandString); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}

	return vlan
}
