package codes

import "fmt"

/*
面向对象
题目1 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
题目2 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，
组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
考察点 ：组合的使用、方法接收者。
*/
// 题目2
type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() string {
	return fmt.Sprintf("EmployeeID: %d, Name: %s, Age: %d", e.EmployeeID, e.Name, e.Age)
}

// 题目1
type Shape interface {
	Area() float32
	Perimeter() float32
}

type ShapeCalc struct {
	area      float32
	perimeter float32
}

type Rectangle struct {
	w int
	h int
	c ShapeCalc
}

type Circle struct {
	r float32
	c ShapeCalc
}

func floatFormat(a float32) float32 {
	str := fmt.Sprintf("%.2f", a)
	var res float32
	fmt.Sscanf(str, "%f", &res)
	return res
}

func (a *Rectangle) Area() float32 {
	a.c.area = float32(a.w * a.h)
	return floatFormat(a.c.area)
}
func (a *Rectangle) Perimeter() float32 {
	a.c.perimeter = float32(2*a.w + 2*a.h)
	return floatFormat(a.c.perimeter)
}

func (a *Circle) Area() float32 {
	a.c.area = 3.14 * a.r * a.r
	return floatFormat(a.c.area)
}
func (a *Circle) Perimeter() float32 {
	a.c.perimeter = 2 * 3.14 * a.r
	return floatFormat(a.c.perimeter)
}

func C3_AreaAndPerimeter(a Shape) ShapeCalc {
	return ShapeCalc{area: a.Area(), perimeter: a.Perimeter()}
}
