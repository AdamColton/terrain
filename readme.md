# Terrain

## Index
* [Requirements](#requirements)

## About
Terrain is meant to build 2-dimensional tile maps. All terrains will have a length, width, slice of the integer ids of the tiles and a 2-dimensional slice of size L x W composed of elements in Vals.

The products of the constructors are meant to be raw materials for simulations or games. For this reason everything has been left public.

For simplicity, the term 'terrain' would probably be 'map' in any other languge. But because map already has a specific meaning in Go, terrain is used to avoid confusion. Never the less, that's exactly what this is - a 2-dimensional tile map.

## Requirements
This package requires github.com/AdamColton/ko.

## func PerlinDiscrete
```go
  func PerlinDiscrete(l,w int, vals []int, smoothing int) (*Terrain)
```
Uses a Perlin process to generate a terrain. Over each smoothing iteration, a tile will change to which ever type the plurality of its neighbors are. If there is a tie, it will choose at random.

## func PerlinContinuous
```go
  func PerlinContinuous(l, w, max int, vals []int, smoothing int) (*Terrain)
```
Uses a Perlin process to generate a terrain. Initially, a 2-dimensional slice is constructed with values in [0:max). Each smoothing iteration, a tile takes on the average value of it's neighbors. After the smoothing iterations are complete, the range [0:max) converted to the range [0:len(vals)). So if max was 100, and len(vals) was 10, a value of 5 would translate to 0 and a vaue of 31 would translate to 3.

Because of this process, values will tend to cluster near the middles values. However if you are using a small number of values (especially if you are only using 2) this will produce output similar to PerlinDiscrete but with slightly larger "blobs".

## func DiamondSquare
```go
  DiamondSquare(l, w, max int) (*Terrain)
```
Uses the diamond-squared algorithm to produce a terrain with values (0,1). This function produces a fairly consistent sized "blobs" relative to the map where the Perlin generators tend to produce fairly consisted sized "blobs" on an absolute scale.

## func Compound
```go
  func Compound(terrain1, terrain2 *Terrain) (*Terrain)
```
Produces a terrain where each tile is the product of the two terrains passed in. This particularly useful is one of the terrains has values (0,1) - such as from a DiamondSquared generator.

## type Terrain
```go
  type Terrain struct{
    L, W int
    Vals []int
    Data [][]int
    Regions []*Region
  }
```