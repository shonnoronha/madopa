![Alt text](https://mdg.imgix.net/assets/images/san-juan-mountains.jpg "a title")

- Parent
  - Grand child

# Hello world

**_hello world *this* is a test for golang parser_**
_test italic_
_test italic 2_

**test bold**
**test bold 2**
**hello _shons_ world**

- Single opening tag
  \*\* Double opening tag

```py
print('hello world from code block')
```

| Command |         Description         |  Example Usage   |
| :-----: | :-------------------------: | :--------------: |
|  `ls`   |   List directory contents   |     `ls -l`      |
|  `cd`   |      Change directory       | `cd /home/user`  |
|  `pwd`  |   Print working directory   |      `pwd`       |
|  `rm`   | Remove files or directories | `rm -rf folder/` |
|  `top`  |  Display running processes  |      `top`       |

test the inline [**Test**](./go.mod) something else random

Does this render some `print('hello world')` code block?

- List 1
- List 2
  - nested 1
  - nested 2
  - nested 3
    - super nested 1
    - super nested 2
      - super super nested 1
- List 3

something else

```go
inputFile := flag.String("input", "", "Input markdown file")
outputFile := flag.String("output", "", "Output HTML file")
serverFlag := flag.Bool("serve", false, "Serve the generated HTML file")
flag.Parse()

if *inputFile == "" {
  fmt.Println("Error: File file is required")
  flag.Usage()
  os.Exit(1)
}

if *outputFile == "" {
  *outputFile = replaceExt(*inputFile, ".html")
}

content, err := os.ReadFile(*inputFile)
if err != nil {
  fmt.Printf("Error while reading File %v\n", err)
  os.Exit(1)
}
```

- Item 1
  - Item 1.1
    - Item 1.1.1
    - Item 1.1.2
  - Item 1.2
- Item 2
  - Item 2.1
    - Item 2.1.1
      - Item 2.1.1.1
