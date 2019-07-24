Ensure approprite setup of postgres database.

For the three <postgres>, add appropriate user, database name on which you are working
	and password to that database.

Check port 12345 is not currently in use.

Check if you have all of the dependencies installed. 
	"encoding/json"
	"net/http"
	"github.com/jinzhu/gorm"
	"strconv"
	"github.com/graphql-go/graphql"	
	_ "github.com/lib/pq"
	"fmt"
If not, install them.

Run the main.go using `go run main.go`