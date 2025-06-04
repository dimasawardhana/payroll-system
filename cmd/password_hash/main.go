package main

import (
	"fmt"
	"payroll-system/internal/utils"
)

func main() {
	// get input password from terminal
	var password string
	fmt.Print("Enter password to hash: ")
	_, err := fmt.Scanln(&password)
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println("Hashed Password:", hashedPassword)
}
