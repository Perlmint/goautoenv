package main

var cmdInit = &Command{
	Usage: `
`,
	Short: "",
	Long:  ``,
	Run:   commandInit,
}

func commandInit(cmd *Command, args []string) {
}
