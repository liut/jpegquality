# Detect JPEG quality from bytes or io.ReadSeeker (like os.File)

Input a JPEG file to detect the quality of the image along with the Quanization Table and DQTs.

[![Build Status](https://travis-ci.org/liut/jpegquality.svg?branch=master)](https://travis-ci.org/liut/jpegquality)

Inpiration from [Estimating Quality](http://fotoforensics.com/tutorial-estq.php).

Fixed bug base on HuangYeWuDeng [ttys3/jpegquality](https://github.com/ttys3/jpegquality/commit/6176ce2bb32baad02c5b3dcd977dbc2eab406312).

## Languages/Dependencies:

Golang

## How to Install/Run:

In order to install and run this program, please follow the steps below:

### 1. Clone the repository to your local device

```bash
git clone https://github.com/liut/jpegquality.git
```

### 2. Install Golang:

```bash
Follow this [link](https://go.dev/doc/install) for help.
```

### 3. Build Executable: 

```bash
cd jpegq && go build
 ```
 
### 4. Run executable and specify input jpeg file.
```bash
jpegq [FILENAME]

Example: jpegq myphoto.jpg                     
```

## How to Contribute:

### 1. Fork the [main repository](https://github.com/liut/jpegquality)

### 2. Clone the repo to your local device

```bash
git clone https://github.com/liut/jpegquality.git
```

### 3. Install Golang:

```bash
Follow this [link](https://go.dev/doc/install) for help.
```

### 5. Create a new branch for your changes.

```bash
git checkout -b <name-of-branch>
```

### 6. Submit a Pull Request to the main repository


## Code usage:

````go

	file, err := os.Open("file.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	j, err := jpegquality.New(file) // or NewWithBytes([]byte)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("jpeg quality %d", j.Quality())

````

## Testing

This program's output accurracy was tested thoroughly in jpegquality_tests.go. Feel free to add anymore tests for edge cases!

## Bugs
Please report any bugs below to ensure contributors can keep this code functional and easy to use


## Command line tool

```sh
go get github.com/liut/jpegquality/cmd/jpegq

jpegq myphoto.jpg

```

### Output example
```
2019/10/03 00:09:08 jpegquality.go:135: Quantization table length 130
2019/10/03 00:09:08 jpegquality.go:144: read bytes 130
2019/10/03 00:09:08 jpegquality.go:161: DQT: table index 0 (luminance), precision: 8
2019/10/03 00:09:08 jpegquality.go:208: tbl 0: 23.90413 572.61780 88.04793
2019/10/03 00:09:08 jpegquality.go:212: aver_quality 88
88
```

## Custom logger for debug
```go

jpegquality.SetLogger(log.New(os.Stderr, "jpegq", log.LstdFlags|log.Lshortfile))

```
