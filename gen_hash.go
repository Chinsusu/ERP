package main; import ("fmt"; "golang.org/x/crypto/bcrypt"); func main() { hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), 12); fmt.Println(string(hash)) }
