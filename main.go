package main

import "fmt"

// Greet 返回个性化问候语
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s! Welcome to Git.", name)
}

func main() {
	fmt.Println(Greet("World"))
	fmt.Println(Greet("world"))
}
