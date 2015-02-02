package gexpect

import (
	"log"
	"os/exec"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//_command Tests
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
