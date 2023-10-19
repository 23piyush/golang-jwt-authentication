package helper

// With import we use "(" and not "{"
import(
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUserType (c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access this resource")
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		// a user can access his own data only
		err = errors.New("Unauthorized to access this resources")
		return err
	}
	err = CheckUserType(c, userType)
    return err

}