package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "snapper"
	app.Usage = "Automate your snapshot testing!"
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "initialize snapper",
			Action: initSnapper,
		},
		{
			Name:   "create",
			Usage:  "create a test",
			Action: create,
		},
		{
			Name:   "test",
			Usage:  "run the tests",
			Action: test,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initSnapper(c *cli.Context) error {
	if _, err := os.Stat("__snapper__"); err != nil {
		if os.IsNotExist(err) {
			// initialize.
			os.Mkdir("__snapper__", os.ModePerm)
			os.Mkdir("__snapper__/tests", os.ModePerm)
		} else {
			fmt.Println("snapper is already initialized")
		}
	} else {
		fmt.Println("snapper is already initialized")
	}
	return nil
}

func create(c *cli.Context) error {
	// We expect a bash input like
	// echo "hello world"
	// curl -X GET api/path
	// python process_files.py
	arg := c.Args().Get(0)
	if arg == "" {
		log.Fatal("Must specify a command")
	}
	fmt.Println("Creating a test with command: ", arg)
	files, err := ioutil.ReadDir("./__snapper__/tests")
	if err != nil {
		log.Fatal(err)
	}
	var id = 0
	for _, f := range files {
		fmt.Println(f.Name())
		i, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}
		if i >= id {
			id = i + 1
		}
	}
	var testDir = "__snapper__/tests/" + strconv.Itoa(id)
	os.Mkdir(testDir, os.ModePerm)
	file, err := os.Create(testDir + "/command.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	fmt.Fprintf(file, arg)
	var commandResult = runCommand(arg)
	fmt.Println("Ran: ", arg, " output:")
	fmt.Println(commandResult)

	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Create", "Exit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	if result == "Create" {
		// Create the first snapshot.
		var snapDir = testDir + "/snaps/"
		os.Mkdir(snapDir, os.ModePerm)
		file, err := os.Create(snapDir + strconv.FormatInt(time.Now().Unix(), 10) + ".txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()
		fmt.Fprintf(file, commandResult)
	}

	fmt.Printf("You choose %q\n", result)

	return nil
}

func runCommand(command string) string {
	// Runs the command and returns the output
	cmd := exec.Command("sh", "-c", command)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(stdoutStderr)
}

func test(c *cli.Context) error {
	// Goes through all the tests and makes sure the outputs are the same.
	files, err := ioutil.ReadDir("./__snapper__/tests")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		runTest(f.Name())
	}
	return nil
}

func runTest(testID string) {

	command, err := ioutil.ReadFile("./__snapper__/tests/" + testID + "/command.txt")
	if err != nil {
		log.Fatal(err)
	}
	var mostRecentSnap = 0
	snapFiles, err := ioutil.ReadDir("./__snapper__/tests/" + testID + "/snaps")
	if err != nil {
		log.Fatal(err)
	}
	for _, snapFile := range snapFiles {
		timestamp, err2 := strconv.Atoi(strings.Split(snapFile.Name(), ".")[0])
		if err2 != nil {
			continue
		}
		if timestamp > mostRecentSnap {
			mostRecentSnap = timestamp
		}
	}

	expectedOutput, err := ioutil.ReadFile("./__snapper__/tests/" + testID + "/snaps/" + strconv.Itoa(mostRecentSnap) + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	var actualOutput = runCommand(string(command))
	fmt.Println("Expected output: " + string(expectedOutput))
	fmt.Println("Actual output: " + string(actualOutput))

	if actualOutput == string(expectedOutput) {
		fmt.Println("passed.")
	} else {
		fmt.Println("failed.")
	}
}
