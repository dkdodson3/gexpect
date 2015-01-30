package gexpect

import (
	"log"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(*testing.T) {
	log.Printf("Testing Hello World... ")
	child, err := Spawn("echo \"Hello World\"")
	if err != nil {
		panic(err)
	}
	err = child.Expect("Hello World")
	if err != nil {
		panic(err)
	}
	log.Printf("Success\n")
}

func TestHelloWorldFailureCase(*testing.T) {
	log.Printf("Testing Hello World Failure case... ")
	child, err := Spawn("echo \"Hello World\"")
	if err != nil {
		panic(err)
	}
	err = child.Expect("YOU WILL NEVER FIND ME")
	if err != nil {
		log.Printf("Success\n")
		return
	}
	panic("Expected an error for TestHelloWorldFailureCase")
}

func TestBiChannel(*testing.T) {
	log.Printf("Testing BiChannel screen... ")
	child, err := Spawn("screen")
	if err != nil {
		panic(err)
	}
	sender, reciever := child.AsyncInteractChannels()
	wait := func(str string) {
		for {
			msg, open := <-reciever
			if !open {
				return
			}
			if strings.Contains(msg, str) {
				return
			}
		}
	}
	sender <- "\n"
	sender <- "echo Hello World\n"
	wait("Hello World")
	sender <- "times\n"
	wait("s")
	sender <- "^D\n"
	log.Printf("Success\n")

}

func TestExpectRegex(*testing.T) {
	log.Printf("Testing ExpectRegex... ")

	child, err := Spawn("/bin/sh times")
	if err != nil {
		panic(err)
	}
	child.ExpectRegex("Name")
	log.Printf("Success\n")

}

func TestCommandStart(*testing.T) {
	log.Printf("Testing Command... ")

	// Doing this allows you to modify the cmd struct prior to execution, for example to add environment variables
	child, err := Command("echo 'Hello World'")
	if err != nil {
		panic(err)
	}
	child.Start()
	child.Expect("Hello World")
	log.Printf("Success\n")
}

func TestExpectFtp(*testing.T) {
	log.Printf("Testing Ftp... ")

	child, err := Spawn("ftp ftp.openbsd.org")
	if err != nil {
		panic(err)
	}
	child.Expect("Name")
	child.SendLine("anonymous")
	child.Expect("Password")
	child.SendLine("pexpect@sourceforge.net")
	child.Expect("ftp> ")
	child.SendLine("cd /pub/OpenBSD/3.7/packages/i386")
	child.Expect("ftp> ")
	child.SendLine("bin")
	child.Expect("ftp> ")
	child.SendLine("prompt")
	child.Expect("ftp> ")
	child.SendLine("pwd")
	child.Expect("ftp> ")
	log.Printf("Success\n")
}

func TestInteractPing(*testing.T) {
	log.Printf("Testing Ping interact... \n")

	child, err := Spawn("ping -c8 8.8.8.8")
	if err != nil {
		panic(err)
	}
	child.Interact()
	log.Printf("Success\n")

}

// Unit Tests

//_command
func Test_spawnBadStringOne(t *testing.T) {
	log.Printf("Testing a bad string")
	_, err := _spawn("don't worry")
	assert.NotNil(t, err)
}

func Test_spawnBadCommandOne(t *testing.T) {
	log.Printf("Testing a bad command")
	_, err := _spawn("foo")
	assert.IsType(t, &exec.Error{}, err)
	assert.Contains(t, err.Error(), "file not found")
}

func Test_spawnBadCommandTwo(t *testing.T) {
	log.Printf("Testing another bad command")
	_, err := _spawn("blah|/4 ")
	assert.IsType(t, &exec.Error{}, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func Test_spawnEmptyArgs(t *testing.T) {
	log.Printf("Testing command with no args")
	_, err := _spawn("")
	assert.Contains(t, err.Error(), "No command given to spawn")
}

func Test_spawnCommand(t *testing.T) {
	log.Printf("Testing a good single command")
	result, _ := _spawn("echo")
	assert.Contains(t, result.Cmd.Path, "/bin/echo")
}

func Test_spawnCommandWithArgs(t *testing.T) {
	log.Printf("Testing a good single command with args")
	result, _ := _spawn("echo \"foo bar\" blah")
	assert.Equal(t, result.Cmd.Path, "/bin/echo")
	assert.Len(t, result.Cmd.Args, 3)
}

func Test_spawnResult(t *testing.T) {
	log.Printf("Testing the return values of a good command")
	result, err := _spawn("echo")
	assert.False(t, structEmpty(result))
	assert.IsType(t, &ExpectSubprocess{}, result)
	assert.Nil(t, err)
}

// Helper Functions

// Functions

func structEmpty(object interface{}) bool {
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	if reflect.ValueOf(object).Kind() == reflect.Struct {
		// Creates an empty copy of the struct object
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}

	return false
}
