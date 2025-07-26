package codes

import (
	"fmt"

	"gorm.io/gorm"
)

/*
## SQL语句练习
### 题目1：基本CRUD操作
- 假设有一个名为 students 的表，包含字段:
 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
  - 要求 ：
    - 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
    - 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
    - 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
    - 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
### 题目2：事务语句
- 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键，
	from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
  - 要求 ：
    - 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/
// 题目1.2
type Account struct {
	gorm.Model
	Name    string
	Balance float32
}
type Transaction struct {
	gorm.Model
	From_account_id uint
	To_account_id   uint
	Amount          float32
}

func acc_create(db *gorm.DB) {
	var accs = []Account{
		{Name: "A1", Balance: 200},
		{Name: "A2", Balance: 100},
		{Name: "A3", Balance: 50},
		{Name: "B", Balance: 0},
	}
	// 创建
	db.AutoMigrate(&Account{}, &Transaction{})

	// 添加数据
	db.Create(accs)
}

func acc_Transaction(db *gorm.DB, At1 string, At2 string, transAmount float32) bool {
	tx := db.Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}

	var ac1 Account
	if err := tx.Where("name=?", At1).First(&ac1).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	if ac1.Balance < transAmount {
		tx.Rollback()
		fmt.Println("账户", At1, "转账", transAmount, " , 余额不足。")
		return false
	}
	var ac2 Account
	if err := tx.Where("name=?", At2).First(&ac2).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	if err := tx.Model(&ac1).Update("balance", ac1.Balance-transAmount).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Model(&ac2).Update("balance", ac2.Balance+transAmount).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	if err := tx.Create(&Transaction{From_account_id: ac1.ID, To_account_id: ac2.ID, Amount: transAmount}).Error; err != nil {
		tx.Rollback()
		fmt.Println("Transaction创建失败。")
		panic(err)
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Println("事务提交失败。")
		return false
	}
	return true
}

func acc_ShowAll(db *gorm.DB) {
	// 账户信息
	var acInfos = []Account{}
	if err := db.Find(&acInfos).Error; err != nil {
		panic(err)
	}
	for _, a := range acInfos {
		fmt.Println("账户信息情况：", a.Name, a.Balance)
	}
}

func Acc_transRun() {
	db := getGormConn()

	// 创建 Account 表
	acc_create(db)

	fmt.Println("转账前账户信息：")
	acc_ShowAll(db)

	var transAccount = []string{"A1", "A2", "A3"}
	trancAmount := float32(100)
	for _, a := range transAccount {
		if acc_Transaction(db, a, "B", trancAmount) {
			fmt.Println("转账执行结果：", a, "转账B, ", trancAmount, " 执行成功")
		} else {
			fmt.Println("转账执行结果：", a, "转账B, ", trancAmount, " 执行失败")
		}
	}
	fmt.Println("转账后账户信息：")
	acc_ShowAll(db)

	// // 删除表中所有数据
	db.Unscoped().Where("1=1").Delete(&Account{})
}

// 题目1.1
type Student struct {
	gorm.Model
	Name  string
	Age   uint8 `gorm:"check:age >= 0 AND age <= 150"`
	Grade string
}

func (a Student) getInfo() map[string]string {
	return map[string]string{"Name": a.Name, "Age": string(a.Age), "Grade": a.Grade}
}

func stu_create(db *gorm.DB) {
	var stus = []Student{
		{Name: "张三", Age: 20, Grade: "三年级"},
		{Name: "李四", Age: 19, Grade: "三年级"},
		{Name: "王五", Age: 14, Grade: "二年级"},
		{Name: "周大", Age: 18, Grade: "一年级"},
		{Name: "小六", Age: 13, Grade: "幼儿园"},
	}
	// 创建
	db.AutoMigrate(&Student{})

	// 添加数据
	db.Create(stus)
}

func Stu_f1() {
	db := getGormConn()

	// 创建 students 表，并添加数据
	stu_create(db)

	var resStus = []Student{}
	// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	IsError(db.Raw("SELECT Name, Age, Grade FROM students WHERE age > ?", 18).Scan(&resStus).Error)
	for _, res := range resStus {
		fmt.Println("age > 18, res:", res.getInfo())
	}

	// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	db.Exec("UPDATE students SET grade=? WHERE name=?", "四年级", "张三")
	IsError(db.Raw("SELECT Name, Grade FROM students WHERE Name = ?", "张三").Scan(&resStus).Error)
	if resStus[0].Grade == "四年级" {
		fmt.Println("张三 的学生年级更新为 四年级 【执行成功】。")
	}

	resStus = []Student{}
	// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	db.Exec("DELETE FROM students WHERE age<?", 15)
	IsError(db.Raw("SELECT * FROM students WHERE age<?", 15).Scan(&resStus).Error)
	if len(resStus) == 0 {
		fmt.Println("删除 students 表中年龄小于 15 岁的学生记录，【执行成功】。")
	}

	// 删除 Student 表
	db.Unscoped().Where("1=1").Delete(&Student{})
}
