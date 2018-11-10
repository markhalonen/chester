package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

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
	fmt.Println("Output:")
	fmt.Println(commandResult)

	prompt := promptui.Select{
		Label: "Create Test?",
		Items: []string{"Create", "Exit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	if result == "Create" {
		// Create the first snapshot.
		file, err := os.Create(testDir + "/expected_output.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()
		fmt.Fprintf(file, commandResult)
		fmt.Println("Test created! Run tests with `snapper test`")
	} else {
		return nil
	}

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

	expectedOutput, err := ioutil.ReadFile("./__snapper__/tests/" + testID + "/expected_output.txt")
	if err != nil {
		log.Fatal(err)
	}
	var actualOutput = runCommand(string(command))

	if actualOutput == string(expectedOutput) {
		fmt.Println(testID, ": passed")
	} else {
		fmt.Println(testID, ": failed")
		printWithBorder("Expected Output", string(expectedOutput))
		printWithBorder("Actual Output", string(actualOutput))
	}
}

func printWithBorder(title, content string) {
	var border = "========================================================="
	title = " " + title + " "
	var titleIdx = (len(border) - len(title)) / 2
	var titleWithBorder = border[0:titleIdx] + title + border[titleIdx+len(title):]
	fmt.Println(titleWithBorder)
	fmt.Println(content)
	fmt.Println(border)
}
