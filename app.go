package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
		log.Fatal("Must specify a command or a directory containing command.sh")
	}

	// First check if it's a folder with command.sh inside it
	// Else run it as a command itself.
	var runDir = ""
	var command = ""
	if _, err := os.Stat(filepath.Join(arg, "command.sh")); !os.IsNotExist(err) {
		// run from here!
		runDir = arg
		command = "./command.sh"
	} else {
		// They passed in a string command directly.
		runDir = ""
		command = arg
	}

	var commandResult = runCommandFromDir(command, runDir)
	printWithBorder("Output", commandResult)

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
		// Create the snapshot.
		createTest(runDir, arg, commandResult)
	} else {
		return nil
	}

	return nil
}

func createTest(runDir, arg, commandResult string) {
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
	testDir := filepath.Join("__snapper__/tests/", strconv.Itoa(id))
	runTestDir := filepath.Join(testDir, "run_test")
	os.MkdirAll(runTestDir, os.ModePerm)

	if runDir == "" {
		// arg should be written to file
		file, err := os.Create(runTestDir + "/command.sh")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()
		fmt.Fprintf(file, arg)
	} else {
		// arg is a folder, it's contents should be copied over
		filesToCopy, _ := filepath.Glob(filepath.Join(arg, "/*"))

		for _, f := range filesToCopy {
			cpCmd := exec.Command("cp", "-r", f, filepath.Join(runTestDir, "/"))

			output, err := cpCmd.CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
				log.Fatal(err)
			}
		}

	}

	file, err := os.Create(testDir + "/expected_output.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	fmt.Fprintf(file, commandResult)
	fmt.Println("Test created! Run tests with `snapper test`")
}

func runCommandFromDir(command, dir string) string {
	// Runs the command and returns the output
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(stdoutStderr)
}

func runCommand(command string) string {
	return runCommandFromDir(command, "")
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
	testDir := "./__snapper__/tests/" + testID
	runTestDir := filepath.Join(testDir, "run_test")
	command, err := ioutil.ReadFile(filepath.Join(runTestDir, "command.sh"))
	if err != nil {
		log.Fatal(err)
	}

	expectedOutput, err := ioutil.ReadFile(testDir + "/expected_output.txt")
	if err != nil {
		log.Fatal(err)
	}
	var actualOutput = runCommandFromDir(string(command), runTestDir)

	if actualOutput == string(expectedOutput) {
		fmt.Println(testID, ": passed")
	} else {
		fmt.Println(testID, ": failed")
		printWithBorder("Expected Output", string(expectedOutput))
		printWithBorder("Actual Output", string(actualOutput))
		prompt := promptui.Select{
			Label: "Options",
			Items: []string{"Update Expected Output", "Delete Test", "Skip", "Exit"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if result == "Exit" {
			os.Exit(0)
		} else if result == "Skip" {
			fmt.Println("Skipping")
			return
		} else if result == "Delete Test" {
			os.RemoveAll(testDir)
		} else if result == "Update Expected Output" {
			ioutil.WriteFile(testDir+"/expected_output.txt", []byte(actualOutput), os.ModePerm)
		}

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
