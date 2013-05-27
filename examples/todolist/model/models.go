package model

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/go-sql-driver/mysql"
    "time"
)

var orm beedb.Model

type User struct {
	Id        int
	Name      string
	Email     string
	Pic       string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateUser(user *User) (err error) {
     err = orm.SetTable("users").Save(user)
     fmt.Println(err)
    return
}


func Get() {
    var user User
    fmt.Println(&user == nil)
    err := orm.SetTable("users").Where(1).Find(&user)
    fmt.Println(err)
    fmt.Println(user.CreatedAt)
}

func Init() {
	db, err := sql.Open("mysql", "root:123456@/todolist?charset=utf8")
	if err != nil {
		panic(err)
	}
    beedb.OnDebug = true
	orm = beedb.New(db)
}

type Task struct {
    Id        int
    UserId   int
    Title      string
    Content    string
    Active    bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

func FindTasks() (tasks []Task) {
    err := orm.SetTable("tasks").Limit(10).Where("user_id>?", 1).FindAll(tasks)
    if err != nil {
        panic(err)
    }
    return
}

func main() {
    //Init()
    user := User{
        Id: 111,
        Name: "我是中国人3",
        Email: "gy@182.com",
        Pic: "agal",
        Active: true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    CreateUser(&user)

    //Get()
}
