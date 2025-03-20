package loader

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "unicode"
  "unicode/utf8"
)

// Create custom split since there some are "," or "." 
// Add "-" for our lovely t-bone
func customSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
  start := 0 
  for start < len(data) {
    r, width := utf8.DecodeRune(data[start:])
    if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
      break
    }
    start += width
  }


  for i := start; i < len(data); {
    r, width := utf8.DecodeRune(data[i:])
    // Check space 
    if unicode.IsSpace(r) {
      return i + width, data[start:i], nil 
    }

    // Add t-bone and Punct 
    if unicode.IsPunct(r) && r != '-' {
      return i + width, data[start:i], nil 
    }

    i += width
  }
  
  if atEOF && start < len(data) {
    return len(data), data[start:], nil 
  }

  return start,nil,nil
}

func TextLoader() (map[string]int32, error) {
  // Open file since it is faster this way 
  file, err := os.Open("food.txt")
  if err != nil {
    return nil, fmt.Errorf("Error opening file: %v", err)
  }

  defer file.Close()

  // Create scanner 
  scanner := bufio.NewScanner(file)
  //scanner.Split(bufio.ScanWords)
  scanner.Split(customSplit)

  targetWords := []string{"t-bone","fatback","pastrami","pork","meatloaf","jowl","enim","bresaola"}

  // Convert list 
  wordSet := make(map[string]struct{})
  for _, word := range targetWords {
    wordSet[word] = struct{}{}
  }

  // Init map word 
  wordCounts := make(map[string]int32)

  // Iterate 
  for scanner.Scan() {
    word := strings.ToLower(scanner.Text()) 
    if _, exists := wordSet[word]; exists {
      wordCounts[word]++
    }
  }

  if err := scanner.Err(); err != nil {
    return nil, fmt.Errorf("Error reading file: %v", err)
  }

//  for _, word := range targetWords {
//    fmt.Printf("%s: %d \n", word, wordCounts[word])
//  }
  return wordCounts, nil
}
