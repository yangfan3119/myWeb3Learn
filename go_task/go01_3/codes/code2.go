package codes

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

/*
Sqlx入门
题目1：使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，
包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，
并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，
并将结果映射到一个 Employee 结构体中。

题目2：实现类型安全映射
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/
// 题目2
type Book struct {
	Id     uint    `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float32 `db:"price"`
}

var cTab_book = `
CREATE TABLE IF NOT EXISTS books (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	price DECIMAL(10,2) NOT NULL
);
`

func book_create(db *sqlx.DB) {
	if _, err := db.Exec(cTab_book); err != nil {
		fmt.Println("创建失败")
		panic(err.Error())
	}

	var bks = []Book{
		{Title: "语文书", Author: "主编1", Price: 32},
		{Title: "数学书", Author: "主编2", Price: 35},
		{Title: "人类简史", Author: "尤瓦尔・赫拉利", Price: 108.0},
		{Title: "西游记", Author: "吴承恩", Price: 79.5},
	}
	if _, err := db.NamedExec(`
	INSERT INTO books (title, author, price)
	VALUES (:title, :author, :price)`, bks); err != nil {
		fmt.Println("导入数据失败")
		panic(err.Error())
	}
}

func Book_Run() {
	db := getSqlxConn()
	// 创建初始化并导入数据
	book_create(db)

	// 查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
	var books = []Book{}
	if err := db.Select(&books, "SELECT * FROM books WHERE price>$1 ORDER BY price DESC", 50); err != nil {
		fmt.Println("查询价格大于50元的书失败。")
		panic(err.Error())
	}
	fmt.Println("查询价格大于50元的书结果信息:", books)

	// 删除所有数据
	if _, err := db.Exec("DELETE FROM books"); err != nil {
		fmt.Println("删除数据失败")
		panic(err.Error())
	}
}

// 题目1
type Employee struct {
	Id         uint    `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float32 `db:"salary"`
}

var cTab_employee = `
CREATE TABLE IF NOT EXISTS employees (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	department TEXT NOT NULL,
	salary DECIMAL(10,2) NOT NULL
);
`

func emp_create(db *sqlx.DB) {

	if _, err := db.Exec(cTab_employee); err != nil {
		fmt.Println("创建失败")
		panic(err.Error())
	}

	var emps = []Employee{
		{Name: "技术1", Department: "技术部", Salary: 12300.5},
		{Name: "技术2", Department: "技术部", Salary: 15100},
		{Name: "销售1", Department: "销售部", Salary: 13300.6},
		{Name: "销售2", Department: "销售部", Salary: 21300.6},
	}
	if _, err := db.NamedExec(`
	INSERT INTO employees (name, department, salary)
	VALUES (:name, :department, :salary)`, emps); err != nil {
		fmt.Println("导入数据失败")
		panic(err.Error())
	}
}

func Emp_run() {
	db := getSqlxConn()
	// 创建初始化并导入数据
	emp_create(db)
	// 使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
	var emps = []Employee{}
	if err := db.Select(&emps, "SELECT * FROM employees WHERE department=$1", "技术部"); err != nil {
		fmt.Println("查询技术部失败。")
		panic(err.Error())
	}
	fmt.Println("查询技术部所有员工信息:", emps)

	var emp = Employee{}
	// 使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	if err := db.Get(&emp, "SELECT * FROM employees WHERE salary=(SELECT MAX(salary) FROM employees)"); err != nil {
		fmt.Println("查询最高工资失败。")
		panic(err.Error())
	}
	fmt.Println("查询工资最高的员工信息:", emp)

	// 删除所有数据
	if _, err := db.Exec("DELETE FROM employees"); err != nil {
		fmt.Println("删除数据失败")
		panic(err.Error())
	}
}
