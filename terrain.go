package terrain

import (
  . "github.com/AdamColton/ko"
  "math/rand"
  "strconv"
  "strings"
)

type Terrain struct{
  L, W int
  Vals []int
  Data [][]int
  Regions []*Region
}

///// Begin Constructors /////

func PerlinDiscrete(l,w int, vals []int, smoothing int) (*Terrain) {
  valsLength := len(vals)
  f_rand := func(i ...int) int{
    return vals[ rand.Intn(valsLength) ]
  }
  self := Terrain{l, w, vals, Slicer(f_rand,l,w).([][]int), make([]*Region, 0)}
  
  f_smooth := func(x,y int) int{
    return self.discreteAvg(x,y)
  }
  
  Looper(smoothing, func() {
    self.Data = Slicer(f_smooth, self.L, self.W).([][]int)
  })

  self.findRegions()

  return &self
}

func PerlinContinuous(l, w, max int, vals []int, smoothing int) (*Terrain) {
  f_rand := func(i ...int) int{
    return rand.Intn(max)
  }
  self := Terrain{l, w, vals, Slicer(f_rand,l,w).([][]int), make([]*Region, 0)}
  
  f_smooth := func(x,y int) int{
    return self.continuousAvg(x,y)
  }
  
  Looper(smoothing, func() {
    self.Data = Slicer(f_smooth, l, w).([][]int)
  })

  conversionFactor := max / len(vals)
  f_convert := func (x,y int) int {
    return vals[ (self.Data[x][y]+1) / conversionFactor ]
  }

  self.Data = Slicer(f_convert, l, w).([][]int)

  self.findRegions()

  return &self 
}

func Compound(terrain1, terrain2 *Terrain) (*Terrain) {
  l, w := terrain1.L, terrain1.W
  if terrain2.L < l { l = terrain2.L }
  if terrain2.W < w { w = terrain2.W }

  funcMult := func(x,y int) (int) {
    return terrain1.Get(x,y) * terrain2.Get(x,y)
  }

  valSet := make(map[int] bool)
  for i := range Product(terrain1.Vals, terrain2.Vals).(chan []int) {
    valSet[ i[0] * i[1] ] = true
  }
  vals := make([]int, 0, len(valSet))
  for i := range valSet {
    vals = append(vals, i)
  }

  self := &Terrain{l, w, vals, Slicer(funcMult, l, w).([][]int), make([]*Region, 0)}
  self.findRegions()
  return self
}

func DiamondSquare(l, w, max int) (*Terrain){
  size := 24
  f_rand := func (i ...int) int {
    return rand.Intn(max)
  }
  data := Slicer(f_rand,size,size).([][]int)

  diamondSquareFunc := func(x, y int) int {
    noise := rand.Intn(max / 3) - (max / 6)
    if x % 2 == 1 {
      if y % 2 == 1 {
        sum := data[(x-1)/2][(y-1)/2]
        sum += data[(x+1)/2][(y-1)/2]
        sum += data[(x-1)/2][(y+1)/2]
        sum += data[(x+1)/2][(y+1)/2]
        return sum / 4
      }
      return noise + ( data[(x-1)/2][y/2] + data[(x+1)/2][y/2] ) / 2
    } else if y % 2 == 1 {
      return noise + ( data[x/2][(y-1)/2] + data[x/2][(y+1)/2] ) / 2
    }
    return noise + data[x/2][y/2]
  }

  for size < l || size < w {
    size += size - 1
    data = Slicer(diamondSquareFunc, size, size).([][]int)
  }

  self := &Terrain{l, w, []int{0,1}, data, make([]*Region,0)}
  avg := 0
  f_smooth := func(x,y int) int{
    r := self.continuousAvg(x,y)
    avg += r
    return r
  }
  self.Data = Slicer(f_smooth, l, w).([][]int)

  max = max/2
  cast := func(x, y int) (int) {
    if self.Data[x][y] > max {
      return 1
    }
    return 0
  }

  self.Data = Slicer(cast, l, w).([][]int)

  self.findRegions()

  return self
}

///// End Constructors /////

func (self *Terrain) String() string {
  stringMap := self.StringMap()
  fn := func(i int) string{
    return strings.Join(stringMap[i], "")
  }
  return strings.Join( Slicer(fn, len(stringMap)).([]string), "\n" ) + "\n"
}

func (self *Terrain) StringMap()([][]string) {
  tileMap := []string{"  ", "..", "__", ";;", ",,", "--", "++", "##", "^^","==", "//", "\\\\"}
  l := len(tileMap) - 1
  fn := func(x,y int) string {
    i := self.Data[x][y]
    if i != 0{
      i = (i % l) + 1
    }
    return tileMap[i]
  }
  return Slicer(fn, self.L, self.W).([][]string)
}

func (self *Terrain) Get(x, y int) (int){
  return self.Data[x % self.L][y % self.W]
}

func (self *Terrain) discreteAvg(x, y int) int {
  counter := make(map[int] int)
  for _, v := range self.Vals {
    counter[v] = 0
  }

  maxVal := 0
  maxIndex := make([]int, 1, 9)

  for i:=-1; i<2; i++ {
    for j:=-1; j<2; j++ {
      val := self.Get(x+i, y+j)
      counter[val]++
      if counter[val] > maxVal {
        maxVal = counter[val]
        maxIndex = maxIndex[:1]
        maxIndex[0] = val
      } else if counter[val] == maxVal {
        maxIndex = append(maxIndex, val)
      }
    }
  }
  return maxIndex[ rand.Intn(len(maxIndex)) ]
}

func (self *Terrain) continuousAvg(x, y int) int {
  sum := 0

  for i:=-1; i<2; i++ {
    for j:=-1; j<2; j++ {
      sum += self.Get(x+i, y+j)
    }
  }
  return sum / 9
}

func (self *Terrain) floodRegion(x,y int, associated map[Coord]bool){
  newRegion := Region{0, make([]Coord, 0)} 
  newRegion.Val = self.Data[x][y]

  dirs := Dirs()
  queue := []Coord{ Coord{x,y} }
  
  for len(queue) > 0 {
    checking := queue[len(queue) - 1]
    queue = queue[:len(queue) - 1]
    newRegion.Coords = append(newRegion.Coords, checking)
    for _, dir := range dirs {
      x, y := checking.x + dir.x, checking.y + dir.y
      if x >= 0 && y >= 0 && x < self.L && y < self.W && self.Data[x][y] == newRegion.Val {
        toCheck := Coord{x,y}
        if _, checked := associated[toCheck]; !checked {
          queue = append(queue, toCheck)
          associated[toCheck] = true
        }
      }
    }
  }
  self.Regions = append(self.Regions, &newRegion)
}

func (self *Terrain) findRegions(){
  associated := make(map[Coord]bool)
  funcCheckAndFlood := func(x,y int) bool{
    if _, checked := associated[Coord{x,y}]; !checked {
      self.floodRegion(x,y, associated)
    }
    return false
  }
  Slicer(funcCheckAndFlood, self.L, self.W)
}

type Region struct{
  Val int
  Coords []Coord
}

func (self *Region) Json(ind string)(string){
  out := "{\n"
  out += ind + "  'val': " + strconv.Itoa(self.Val) + ",\n"
  out += ind + "  'coords': [\n"

  getCoords := func (i int) string {
    return ind + "    [" + self.Coords[i].String() + "]"
  }
  coords := Slicer(getCoords, len(self.Coords)).([]string)
  out += strings.Join(coords, ",\n")

  out += "\n" + ind +"  ]\n"
  return out + ind + "}"
}

type Coord struct{
  x,y int
}

func NewCoord(x,y int) Coord{
  return Coord{x,y}
}

func (self *Coord) X()int{ return self.x }
func (self *Coord) Y()int{ return self.y }
func (self *Coord) String () (string){
  return strconv.Itoa(self.x) + ", " +strconv.Itoa(self.y)
}

func Dirs() []*Coord {
  return []*Coord{ 
    &Coord{-1,0},
    &Coord{1,0},
    &Coord{0,-1},
    &Coord{0,1}}
}