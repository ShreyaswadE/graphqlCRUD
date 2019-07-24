package main

import (
	"encoding/json"
	"net/http"
	"github.com/jinzhu/gorm"
	"strconv"
	"github.com/graphql-go/graphql"	
	_ "github.com/lib/pq"
	"fmt"
)

type Employee struct{
  EmpNo int `json:"EmpNo" gorm:"AUTO_INCREMENT" sql:"unique;not null" gorm:"primary_key"`
  Ename string `json:"ename" form:"eName"`
  Job string `json:"job" form:"job"`
  Mgr string `json:"mgr" form:"mgr"`
  Salary int `json:"salary" form:"salary"`
  DeptNo int `json:"deptno"`
}

type Department struct{
	DeptNo int `json:"no" gorm:"primary_key" sql:"unique;not null"`
	Dname string `json:"dname"`
	Loc  string  `json:"location"`
	Employees []Employee
}

func createData(db *gorm.DB) {
	var emp[] Employee = []Employee {
		Employee{
			Ename : "shreyas",
			Job : "computer eng",
			Mgr : "hello",
			Salary : 788,
			DeptNo : 1,
		},
		Employee{
			Ename : "wade",
			Job : "entc ent",
			Mgr : "heaven",
			Salary : 3,
			DeptNo : 2,
		},
		Employee{
			Ename : "mukta",
			Job : "barve",
			Mgr : "hellp planntet",
			Salary : 3,
			DeptNo : 1,
		},
			Employee{
			Ename : "ankit",
			Job : "gupta",
			Mgr : "wolrd hello",
			Salary : 3,
			DeptNo : 2,
		},
		Employee{
			Ename : "naina",
			Job : "da kasoor",
			Mgr : "andhadoon",
			Salary : 3,
			DeptNo : 1,
		},
		Employee{
			Ename : "kabir",
			Job : "singh",
			Mgr : "hmm wolrd",
			Salary : 3,
			DeptNo : 2,
		},
	}

	var dept[] Department = []Department {
		Department{
			DeptNo : 1,
			Dname : "comp",
			Loc : "Pune",
		},
		Department{
			DeptNo : 2,
			Dname : "entc",
			Loc : "Pune",
		},
	}
	if err := db.Debug().Create(&dept[0]); err.Error != nil{
		panic(err.Error)
	}
	if err := db.Debug().Create(&dept[1]); err.Error!= nil{
		panic(err.Error)
	}
	for _ , employee := range emp {
		if err := db.Debug().Create(&employee); err.Error != nil{
			panic(err.Error)
		}
	}
} 


func initTables() *gorm.DB {
	db, err := gorm.Open("postgres", "user=<postgres> dbname=<postgres> password=<postgres> sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("database open ")

	employee := Employee{}
	department := Department{}
	if err := db.Debug().DropTable(&employee); err.Error != nil{
		panic(err.Error)
	}
	if err := db.Debug().CreateTable(&employee); err.Error != nil{
		panic(err.Error)
	}
	if err := db.Debug().DropTable(&department); err.Error != nil{
		panic(err.Error)
	}
	if err := db.Debug().CreateTable(&department); err.Error != nil{
		panic(err.Error)
	}

	createData(db)
	return db
}

func resolveAllEmployees(db *gorm.DB, params graphql.ResolveParams) (interface{}, error){
	employees := []Employee{}
	db.Debug().Find(&employees)
    if err:= db.Debug().Find(&employees); err.Error!= nil {
        return nil,err.Error
   	}
    fmt.Println(employees)
    return employees,nil
}

func getEmployeeFromId(db *gorm.DB, params graphql.ResolveParams) (interface{}, error){
	employeeFound := Employee{}
	idValue,_ := strconv.Atoi(params.Args["id"].(string))
    if err :=db.Where("emp_no = ?", idValue).First(&employeeFound); err.Error!= nil {
        return nil,err.Error
    }
    fmt.Println(employeeFound)
    return employeeFound,nil
}

func getAllEmployeesInDept (db *gorm.DB, params graphql.ResolveParams)(interface{}, error){
	employees := []Employee{}
	if err := db.Debug().Model(&Employee{}).Joins("inner join departments on employees.dept_no = departments.dept_no").
					Where("departments.dname = ?",params.Args["deptname"].(string)).
	        		Select("employees.emp_no,employees.ename,employees.job,employees.mgr,employees.salary,employees.dept_no").
	        		Scan(&employees); err.Error != nil{
	        			return nil,err.Error
	    		}

    return employees,nil
}

func createNewEmployee(db *gorm.DB, params graphql.ResolveParams) (interface{},error){
	ename, _ := params.Args["ename"].(string)
	job, _ := params.Args["job"].(string)
	mgr, _ := params.Args["mgr"].(string)
	salary, _ := strconv.Atoi(params.Args["salary"].(string))
	deptno, _ :=  strconv.Atoi(params.Args["deptno"].(string))

	var employee = Employee{
		Ename:ename,
		Job: job,
		Mgr:mgr,
		Salary:salary,
		DeptNo: deptno,
	}
	
	if err:=db.Debug().Create(&employee); err.Error != nil{
		return nil,err.Error
	}
	return employee,nil
}

func removeEmployeeById(db *gorm.DB, params graphql.ResolveParams) (interface{},error){
	id, _ := strconv.Atoi(params.Args["id"].(string))
	
	if err := db.Debug().Where("emp_no = ?", id).Delete(&Employee{}); err.Error != nil {
		return nil,err.Error
	}
	employees := []Employee{}
   	
	if err := db.Debug().Find(&employees); err.Error != nil {
	    		return nil,err.Error
         		}
       		return employees,nil
}

func main() {

	db := initTables()

	db.Debug().Model(&Employee{}).AddForeignKey("dept_no","departments(dept_no)","CASCADE","CASCADE")

	employeeType := graphql.NewObject(graphql.ObjectConfig{
	Name: "Employee",
	Fields: graphql.Fields{
			"empno": &graphql.Field{
				Type: graphql.String,
			},
			"ename": &graphql.Field{
				Type: graphql.String,
			},
			"job": &graphql.Field{
				Type: graphql.String,
			},
			"mgr": &graphql.Field{
				Type: graphql.String,
			},
			"salary": &graphql.Field{
				Type: graphql.String,
			},
			"deptNo": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"employees": &graphql.Field{
    	    Type: graphql.NewList(employeeType),
	        Resolve: func(params graphql.ResolveParams) (interface{}, error) {
	        	return resolveAllEmployees(db, params)
        	},
    	},
    	"employee": &graphql.Field{
    	    Type: employeeType,
    	    Args: graphql.FieldConfigArgument{
            	"id": &graphql.ArgumentConfig{
           	     	Type: graphql.String,
          	  	},
        	},
	        Resolve: func(params graphql.ResolveParams) (interface{}, error) {
	        	return getEmployeeFromId(db , params)
        	},
    	},
    	"deptEmployee": &graphql.Field{
    	    Type: graphql.NewList(employeeType),
    	    Args: graphql.FieldConfigArgument{
            	"deptname": &graphql.ArgumentConfig{
           	     	Type: graphql.String,
          	  	},
        	},
	        Resolve: func(params graphql.ResolveParams) (interface{}, error) {
	        	return getAllEmployeesInDept(db,params)
        	},
    	},
    }})

    rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"enterEmployee": &graphql.Field{
				Type: employeeType, // the return type for this field
				Args: graphql.FieldConfigArgument{
					"ename": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"job": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"mgr": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"salary": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"deptno": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{},error) {
					return createNewEmployee(db, params)
				},
			},
			"removeEmployee": &graphql.Field{
				Type:  graphql.NewList(employeeType), // the return type for this field
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{},error) {
					return removeEmployeeById(db,params)
				},
			},
		},
	})


	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}) 	

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":12345", nil)


}