package main

import (
	"fmt"
	// orm  Object Relational Mapping 对象关系映射
	"github.com/jinzhu/gorm"
	// mysql驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// UserInfo 用户信息
type UserInfo struct {
	ID     uint
	Name   string
	Gender string
	Hobby  string
}

// todo 测试
func add(db *gorm.DB, name string) {
	// Create a new TestUser
	u1 := UserInfo{1, name, "男", "篮球"}
	u2 := UserInfo{2, "沙河娜扎", "女", "足球"}
	// 创建记录
	db.Create(&u1)
	db.Create(&u2)
}

// todo 测试

func remove(db *gorm.DB, id uint) {
	// Delete content based on ID
	u := new(UserInfo)
	u.ID = id
	db.First(u)
	db.Delete(&u)
	fmt.Println("删除成功")
}

// todo 测试
func update(db *gorm.DB, id uint, newName string) bool {
	u := new(UserInfo)
	u.ID = id
	db.Model(&u).Update("Name", newName)
	return true
}

// todo 测试
func select_user(db *gorm.DB, id uint) *UserInfo {
	var u = new(UserInfo)
	u.ID = id
	db.First(u)
	return u
}

func main() {

	// connetion mysql inif
	user := "root"
	password := "root1234"
	host := "172.28.233.113"
	port := "13306"
	dbname := "testdb"
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	// 连接数据库
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// 结束前关闭
	defer db.Close()

	// 自动迁移
	db.AutoMigrate(&UserInfo{})

	// todo 测试
	// 查
	u := select_user(db, 1)
	fmt.Printf("u:%#v\n", u)

	// 增
	fmt.Println("添加 id=1")
	add(db, "旧名字")

	// 改
	fmt.Println("修改 id=1")
	update(db, 1, "新名字")

	// 删
	fmt.Println("删除 id=2")
	remove(db, 2)

	// 查
	uu := select_user(db, 1)
	fmt.Printf("uu:%#v\n", uu)

	// // 查询
	// var u = new(UserInfo)
	// db.First(u)
	// fmt.Printf("%#v\n", u)

	// var uu UserInfo
	// db.Find(&uu, "hobby=?", "足球")
	// fmt.Printf("%#v\n", uu)

	// // 更新
	// // 删除
	// db.Delete(&u)
}
