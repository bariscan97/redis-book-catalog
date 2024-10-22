package middleware

// import (
// 	// "strings"
// 	"fmt"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// func SearchFields(c *gin.Context) error {
	
// 	query := ""

// 	queries := c.Request.URL.Query()

// 	for key ,value := range queries {
		
// 		if key == "page" {
// 			continue
// 		}

// 		switch key{
// 		case "gt":
// 			_ , err := strconv.Atoi(value[0])
// 			if err != nil {
// 				return fmt.Errorf(err.Error())
// 			}
// 			query += fmt.Sprintf("@%s:[%d +inf] ", key, value)
	
// 		case "lt":
// 			_ , err := strconv.Atoi(value[0])
// 			if err != nil {
// 				return fmt.Errorf(err.Error())
// 			}
// 			query += fmt.Sprintf("@%s:[0 %d] ", key, value)
		
// 		default:
// 			query += fmt.Sprintf("@%s:%s ", key, value)
// 		}
// 	}

// 	if query == "" {
// 		query = "*"
// 	}

// 	c.Next()
// }