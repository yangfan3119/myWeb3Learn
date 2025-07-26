package codes

import "testing"

func Test_PrintInfo(t *testing.T) {
	type test struct {
		input  Employee
		output string
	}
	tests := map[string]test{
		"No1": {input: Employee{Person: Person{Name: "张三", Age: 18}, EmployeeID: 1}, output: string("EmployeeID: 1, Name: 张三, Age: 18")},
		"No2": {input: Employee{Person: Person{Name: "李四"}, EmployeeID: 23}, output: string("EmployeeID: 23, Name: 李四, Age: 0")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.PrintInfo()
			if got != tc.output {
				t.Errorf("exceted: %s; got: %s", tc.output, got)
			}
		})
	}
}

func Test_C3_AreaAndPerimeter(t *testing.T) {
	type test struct {
		input  Shape
		output ShapeCalc
	}
	isEqual := func(a ShapeCalc, b ShapeCalc) bool {
		if (a.area == b.area) && (a.perimeter == b.perimeter) {
			return true
		}
		return false
	}

	tests := map[string]test{
		"No1": {input: &(Rectangle{w: 2, h: 3}), output: ShapeCalc{area: 6, perimeter: 10}},
		"No2": {input: &(Rectangle{w: 10, h: 10}), output: ShapeCalc{area: 100, perimeter: 40}},
		"No3": {input: &(Circle{r: 4}), output: ShapeCalc{area: 50.24, perimeter: 25.12}},
		"No4": {input: &(Circle{r: 7}), output: ShapeCalc{area: 153.86, perimeter: 43.96}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := C3_AreaAndPerimeter(tc.input)
			if !isEqual(got, tc.output) {
				t.Errorf("exceted area: %f, peri:%f   got area: %f, peri:%f", tc.output.area, tc.output.perimeter, got.area, got.perimeter)
			}
		})
	}
}
