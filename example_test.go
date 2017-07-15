package env2flag_test

import (
	"flag"
	"fmt"
	"os"

	env2flag "github.com/koron/go-env2flag"
)

func ExampleParse() {
	name := flag.String("my_name", "John Doe", "please set name")
	msg := flag.String("message", "Have a good day.", "please set message")

	os.Setenv("MY_NAME", "George")
	env2flag.Parse()

	fmt.Printf("Hello %s!\n", *name)
	fmt.Println(*msg)
	// Output:
	// Hello George!
	// Have a good day.

}
