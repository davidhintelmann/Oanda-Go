package connect

// Be careful not to expose your password to the public
// modify password_edit.go file found in connect directory
// rename that file password.go and this file only has one function.
// This functions name needs to be edited by removing the underscore
// at the end of 'ImportPassword_' function.
// modify the return statement for the password of your PostgreSQL
// user that is accessing your local instance
func ImportPassword_() string { // delete underscore
	return "XXXXXXXXX" // replace with PostgreSQL password
}
