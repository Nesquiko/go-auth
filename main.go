package main

import "github.com/Nesquiko/go-auth/pkg/app"

func main() {
	// password := []byte("123")

	// // Hashing the password with the default cost of 10
	// hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(hashedPassword))
	// fmt.Printf("len is %d\n", len(hashedPassword))

	// // Comparing the password with the hash
	// err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	// fmt.Println(err) // nil means it is a match
	app.StartServer()
}
