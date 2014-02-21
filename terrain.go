package terrain

import (
  . "github.com/AdamColton/ko"
  "math/rand"
  "strconv"
  "strings"
)

type Terrain struct{
  l,w int
  vals []int
  Data [][]int
  Regions []*Region
}

func (self *Terrain) Size()(int, int){
  return self.l, self.w
}

func (self *Terrain) Vals()([]int){
  ret := make([]int, len(self.vals))
  copy(ret, self.vals)
  return ret
}

func New(l,w int, vals []int, smoothing int) *Terrain{
  valsLength := len(vals)
  f_rand := func(i ...int) int{
    return vals[ rand.Intn(valsLength) ]
  }
  self := Terrain{l, w, vals, Slicer(f_rand,l,w).([][]int), make([]*Region, 0)}
  
  f_smooth := func(x,y int) int{
    return self.stochaticAvg(x,y)
  }
  
  Looper(smoothing, func() {
    self.Data = Slicer(f_smooth, self.l, self.w).([][]int)
  })

  self.findRegions()

  return &self
}

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
  return Slicer(fn, self.l, self.w).([][]string)
}

func (self *Terrain) stochaticAvg(x,y int) int {
  counter := make(map[int] int)
  for _,v := range self.vals{
    counter[v] = 0
  }

  maxVal := 0
  maxIndex := 0

  for i:=-1; i<2; i++ {
    for j:=-1; j<2; j++ {
      if val, ok := self.Get(x+i, y+j); ok{
        counter[val]++
        if counter[val] > maxVal {
          maxVal = counter[val]
          maxIndex = val
        }
      }
    }
  }
  return maxIndex
}

func (self *Terrain) Get(l,w int) (int, bool){
  if (l > 0 && l < self.l && w > 0 && w < self.w){
    return self.Data[l][w], true
  }
  return 0, false
}

func Compound(map1, map2 *Terrain) (*Terrain){
  funcMult := func(x,y int) int{
    v1,_ := map1.Get(x,y)
    v2,_ := map2.Get(x,y)
    return v1*v2 
  }

  valSet := make(map[int] bool)
  for i := range Product(map1.vals, map2.vals).(chan []int) {
    valSet[ i[0] * i[1] ] = true
  }

  vals := make([]int, 0, len(valSet))
  for i := range valSet {
    vals = append(vals, i)
  }

  self := &Terrain{map1.l, map1.w, vals, Slicer(funcMult, map1.l, map1.w).([][]int), make([]*Region, 0)}
  self.findRegions()
  return self
}

type Region struct{
  Val int
  Coords []Coord
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

func (self *Terrain) findRegions(){
  associated := make(map[Coord]bool)
  funcCheckAndFlood := func(x,y int) bool{
    if _, checked := associated[Coord{x,y}]; !checked {
      self.floodRegion(x,y, associated)
    }
    return false
  }
  Slicer(funcCheckAndFlood, self.l, self.w)
}

func Dirs() []*Coord {
  return []*Coord{
    &Coord{-1,0},
    &Coord{1,0},
    &Coord{0,-1},
    &Coord{0,1}}
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
      if x>=0 && y>=0 && x<self.l && y<self.w && self.Data[x][y] == newRegion.Val{
        toCheck := Coord{x,y}
        if _, checked := associated[toCheck]; !checked{
          queue = append(queue, toCheck)
          associated[toCheck] = true
        }
      }
    }
  }
  self.Regions = append(self.Regions, &newRegion)
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

func Perlin(l, w, max int, vals []int, smoothing int)(*Terrain){
  f_rand := func(i ...int) int{
    return rand.Intn(max)
  }
  self := Terrain{l, w, vals, Slicer(f_rand,l,w).([][]int), make([]*Region, 0)}
  
  f_smooth := func(x,y int) int{
    return self.perlinAvg(x,y)
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

func (self *Terrain) perlinAvg(x,y int) int {
  sum := 0
  c := 0

  for i:=-1; i<2; i++ {
    for j:=-1; j<2; j++ {
      if val, ok := self.Get(x+i, y+j); ok{
        sum += val
        c++
      }
    }
  }
  return sum / c
}

func DiamondSquare(l, w, max int)(*Terrain){
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
    r := self.perlinAvg(x,y)
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