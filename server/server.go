package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"

	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

func debugLabel(label string, words string) {
	fmt.Println("["+label+"]", words)
}

func systemLogs(words string) {
	debugLabel("SYSTEM", words)
}

type User struct {
	Name  string `json:"name"`
	Value int    `json:"value" default:"-1"`
	Id    int    `json:"id" default:"-1"`
}

type ResponseUser struct {
	Status  bool   `json:"status" default:"false"`
	Message string `json:"msg"`
	Data    User   `json:"data"`
}

type ResponseUsers struct {
	Status  bool   `json:"status" default:"false"`
	Message string `json:"msg"`
	Data    []User `json:"data"`
}

func badRequestUser(re render.Render, message string) {
	res := ResponseUser{Status: false, Message: message}
	re.JSON(http.StatusBadRequest, res)
}

func getLastId(users []User) int {
	maxId := 0
	for i := 0; i < len(users); i++ {
		if users[i].Id > maxId {
			maxId = users[i].Id
		}
	}

	return maxId
}

func writeNewUsers(allUsers []User) error {
	filePath := getFilePath()
	jsonString, _ := json.Marshal(allUsers)
	return ioutil.WriteFile(filePath, jsonString, os.ModePerm)
}

func addNewUsers(newUser User) error {
	allUsers, _ := getAllUsers()
	allUsers = append(allUsers, newUser)
	fmt.Println("allUsers => ", allUsers)

	return writeNewUsers(allUsers)
}

func getAllUsersData() []User {
	allUsers, status := getAllUsers()
	if !status {
		allUsers = []User{}
	}

	return allUsers
}

func updateUser(user User) error {
	allUsers := getAllUsersData()
	for i := 0; i < len(allUsers); i++ {
		if allUsers[i].Id == user.Id {
			allUsers[i] = user
		}
	}

	// write file
	return writeNewUsers(allUsers)
}

func martiniRoute(m *martini.ClassicMartini) {

	m.Group("/upload", func(r martini.Router) {
		r.Put("/file", func(re render.Render) {
			re.Error(http.StatusMethodNotAllowed)
		})
	})

	m.Group("/user", func(r martini.Router) {
		// Create new user
		r.Put("", binding.Bind(User{}), func(user User, re render.Render) {
			allUsers := getAllUsersData()

			// get user info
			// find the last user id(ignore id from input request)
			user.Id = getLastId(allUsers) + 1
			fmt.Println("user => ", user)

			if err := addNewUsers(user); err != nil {
				re.JSON(http.StatusInternalServerError, ResponseUser{Status: false, Message: err.Error()})
			} else {
				re.JSON(http.StatusOK, ResponseUser{Status: true, Data: user})
			}
		})

		r.Delete("/:id", func(re render.Render, params martini.Params) {
			allUsers := getAllUsersData()
			idStr := params["id"]
			if id, err := strconv.Atoi(idStr); err != nil {
				badRequestUser(re, "Cannot found id: "+idStr)
			} else {
				systemLogs("ID:::" + strconv.Itoa(id))
				// try to get user by id
				if user, status := getUserById(id, allUsers); status {
					// remove
					if err := removeUserById(id, allUsers); err == nil {
						re.JSON(http.StatusOK, ResponseUser{
							Message: "Delete success",
							Status:  true,
							Data:    user,
						})
					} else {
						re.JSON(http.StatusInternalServerError, ResponseUser{Status: false, Message: err.Error()})
					}
				} else {
					badRequestUser(re, "Cannot found id: "+idStr)
				}
			}
		})

		// Update user data
		r.Post("/:id", binding.Bind(User{}), func(user User, re render.Render, params martini.Params) {
			fmt.Println("user", user)
			idStr := params["id"]
			if id, err := strconv.Atoi(idStr); err == nil {
				allUsers := getAllUsersData()
				if userOrg, status := getUserById(id, allUsers); status {
					fmt.Println("userOrg", userOrg)
					if user.Name != "" {
						userOrg.Name = user.Name
					}

					if user.Value != 0 {
						userOrg.Value = user.Value
					}

					if err := updateUser(userOrg); err != nil {
						re.JSON(http.StatusInternalServerError, ResponseUser{Status: false, Message: "Try to create a new user"})
					} else {
						re.JSON(http.StatusOK, ResponseUser{Status: true, Data: userOrg})
					}
				} else {
					badRequestUser(re, "Cannot found id: "+idStr)
				}
			} else {
				badRequestUser(re, "ID must be number")
			}
		})

		r.Get("/:id", func(re render.Render, params martini.Params) {
			allUsers := getAllUsersData()
			idStr := params["id"]
			if id, err := strconv.Atoi(idStr); err != nil {
				badRequestUser(re, "Cannot found id: "+idStr)
			} else {
				systemLogs("ID:::" + strconv.Itoa(id))
				if user, status := getUserById(id, allUsers); status {
					re.JSON(http.StatusOK, user)
				} else {
					badRequestUser(re, "Cannot found id: "+idStr)
				}
			}
		})
	})

	m.Get("/users", func(r render.Render) {
		if users, status := getAllUsers(); status {
			res := ResponseUsers{Status: true, Data: users}
			r.JSON(http.StatusOK, res)
		} else {
			r.Error(http.StatusInternalServerError)
		}
	})

	// m.Get("/:name", func(params martini.Params) string {
	// 	return "Hello " + params["name"]
	// })

	m.Any("/", func() string {
		return "Hello world! => NaJa"
	})
}

func getUserById(id int, allUsers []User) (User, bool) {
	data := User{}
	status := false
	for i := 0; i < len(allUsers); i++ {
		if allUsers[i].Id == id {
			data = allUsers[i]
			status = true
			break
		}
	}

	return data, status
}

func removeUserById(id int, allUsers []User) error {
	index := 0

	// find index
	for i := 0; i < len(allUsers); i++ {
		if allUsers[i].Id == id {
			index = i
			break
		}
	}

	// remove array item by index
	newUsers := append(allUsers[:index], allUsers[index+1:]...)
	fmt.Println("newUsers", newUsers)

	return writeNewUsers(newUsers)
}

func getAllUsers() ([]User, bool) {
	dataUsers := []User{}
	status := false
	if dirName, err := os.Getwd(); err == nil {
		status = true
		filePath := dirName + `\assets\users.json`
		var allUsers = "[]"
		if usersFiles, err := ioutil.ReadFile(filePath); err == nil {
			allUsers = string(usersFiles)
		}

		json.Unmarshal([]byte(allUsers), &dataUsers)
	}

	return dataUsers, status
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getFilePath() string {
	dirName, _ := os.Getwd()
	return dirName + `\assets\users.json`
}

func initFile() {
	filePath := getFilePath()
	systemLogs("File Path:" + filePath)
	if usersFiles, err := ioutil.ReadFile(filePath); err != nil {
		// init file

		// create directory
		dirName, _ := os.Getwd()
		dirPath := dirName + `\assets\`
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			if errCreateDir := os.Mkdir(dirPath, os.ModeDir); errCreateDir == nil {
				// create a new file
				var str []byte = []byte("[]")
				if err := ioutil.WriteFile(filePath, str, os.ModePerm); err != nil {
					systemLogs("Cannot create file")
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
	} else {
		// can read file
		systemLogs("Users File::" + string(usersFiles))
	}
}

func main() {
	// init class
	m := martini.Classic()
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"}, // allow all hosts
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	m.Use(render.Renderer())
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	// m.Use(martini.Static("assets")) // serve from the 'assets' directory
	initFile()
	systemLogs("Server has initialized")

	// route
	martiniRoute(m)

	systemLogs("Before running server")
	m.Run()
}
