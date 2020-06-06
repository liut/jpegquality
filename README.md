# Detect JPEG quality from bytes or io.ReadSeeker (like os.File)

[![Build Status](https://travis-ci.org/liut/jpegquality.svg?branch=master)](https://travis-ci.org/liut/jpegquality)

Inpiration from [Estimating Quality](http://fotoforensics.com/tutorial-estq.php).

Fixed bug base on HuangYeWuDeng [ttys3/jpegquality](https://github.com/ttys3/jpegquality/commit/6176ce2bb32baad02c5b3dcd977dbc2eab406312).


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
