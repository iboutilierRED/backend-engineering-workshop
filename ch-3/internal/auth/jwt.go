package auth

const superSecret = "mechaleckahi"

func CreateToken() {

}

// func ValidateToken(tokenStr string) error {
// 	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		// check to see if correct signing method used
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, errors.New(("Unexpected signing method"))
// 		}

// 		return []byte(superSecret), nil
// 	})

// 	if err != nil {
// 		return
// 	}
// }
