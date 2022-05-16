package main

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

type StructTreeMap struct {
	data *treemap.Map
}

func (treemap_ *StructTreeMap) Init() {
	treemap_.data = treemap.NewWith(utils.Float64Comparator)
}

func TestTreeMap() {
	data := treemap.NewWith(utils.Float64Comparator)

	var i = float64(0)
	for i = float64(0); i < 10; i++ {
		data.Put(i*1.1, "--data-"+fmt.Sprint(i))
	}

	iter := data.Iterator()

	// 迭代访问
	fmt.Println("\n顺序访问")
	for iter.Begin(); iter.Next(); {
		fmt.Printf("key: %v, Value: %v \n", iter.Key(), iter.Value())
	}

	fmt.Println("\n逆序访问!")
	for iter.End(); iter.Prev(); {
		fmt.Printf("key: %v, Value: %v \n", iter.Key(), iter.Value())
	}

	// 头尾访问
	iter.First()

	fmt.Println(iter.Key())

	// if iter.Next() {
	// 	fmt.Println(iter.Key())
	// }

	// iter.Last()

	// fmt.Println(iter.Key())

	// if iter.Next() {
	// 	fmt.Println(iter.Key())
	// }

	// fmt.Println((data))
}

// func main() {
// 	fmt.Println("Test Risk Ctrl")

// 	// aggregate.TestInnerDepth()

// 	// aggregate.TestImport()

// 	TestTreeMap()
// }
