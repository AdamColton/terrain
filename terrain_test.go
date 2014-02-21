package terrain

import(
  "fmt"
)

func ExamplePerlinDiscrete() {
  perlinDiscrete := PerlinDiscrete(100, 100, []int{0,1,2,3,4,5}, 4)
  fmt.Println(perlinDiscrete)
}

func ExamplePerlinContinuous() {
  perlinContinuous := PerlinContinuous(100,100, 100000, []int{0,1,2,3,4,5}, 5)
  fmt.Println(perlinContinuous)
}

func ExampleDiamondSquare () {
  diamondSquare := DiamondSquare(100, 100, 100000)
  fmt.Println(diamondSquare)
}

func ExampleCompound() {
  perlinDiscrete := PerlinDiscrete(100, 100, []int{0,1,2,3,4,5}, 4)
  diamondSquare := DiamondSquare(100, 100, 100000)
  compound := Compound(perlinDiscrete, diamondSquare)
  fmt.Println(compound)
}