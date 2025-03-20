package main 

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "log"
  "os"
 // "strings"
)

// Exercise 1: Find the biggest path 
// Implementation: Find the path from bottom to up and keep swap the biggest value
// Note: I am an idiot who misunderstood the assignment that I try to implement BFS search in the first place ._.
func maxPath(triangle [][]int) int {
  // We start from the second last so that we can compare to last row
  for i := len(triangle) -2; i >= 0; i-- {
    for j := 0; j < len(triangle[i]); j++ {
      // Compare left to right node and sum up the parent node with the most value
      if triangle[i+1][j] > triangle[i+1][j+1] {
        triangle[i][j] += triangle[i+1][j]
      } else {
        triangle[i][j] += triangle[i+1][j+1]
      }

    }
  }

  return triangle[0][0]
}

// This is for first test 
// Implementation: Read the file "hard.json" 
// Assign the value I use triangle since node is form in shape of triangle
func testOne() {
  // open json file 
  file, err := os.Open("hard.json")
  if err!= nil {
    log.Fatalf("Failed to open file: %s", err)
  }

  defer file.Close()

  value, err := ioutil.ReadAll(file)
  if err != nil {
    log.Fatalf("Failed to read file: %s", err)
  }

  var triangle [][]int
  err = json.Unmarshal(value, &triangle) 
  if err != nil {
    log.Fatalf("Failed to parse JSON: %s \n", err)
  }

  result := maxPath(triangle)
  fmt.Printf("The maximum path is: %d \n", result)
}

func main() {
  testOne()
}
