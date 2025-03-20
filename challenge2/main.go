package main 

import (
  "fmt"
)

// Test two
func translate_seq(pattern string) string {
  //fmt.Println("Pattern: ", pattern)
  totalLength := len(pattern)
  result := make([]int, totalLength + 1)
 
  //fmt.Println("Init Result: ", result)

  // Check from right -> left and find L first 
  for i := totalLength - 1; i >= 0; i-- {
    switch pattern[i] {
      case 'L':
        if result[i] <= result[i+1] {
        result[i] = result[i+1] + 1 
      }
      case '=': 
        result[i] = result[i+1]
      }
    // No action
  }

  // Check from left -> right and find R 
  //fmt.Println("Left Check Result: ", result)

  for i := 0; i < totalLength; i++ {
    switch pattern[i] {
    case 'R': 
      if result[i] >= result[i+1] {
        result[i+1] = result[i] + 1
      }
    case '=': 
      result[i+1] = result[i]
    }
  }

  //fmt.Println("Right check result: ", result)

  output := make([]byte, totalLength+1) 
  for i, val := range result {
    output[i] = byte('0' + val)
  }

  return string(output)
}

func main() {
  var input string 
  fmt.Print("Enter sequence (only L, R, = allowed): ")
  fmt.Scanln(&input)

  // Validate input 
  validInput := ""
  for _, char := range input {
    if char == 'L' || char == 'R' || char == '=' {
      validInput += string(char)
    }
  }

  res := translate_seq(validInput)
  fmt.Println("Sequence is ", validInput)
  fmt.Println("Output is ", res)
}
