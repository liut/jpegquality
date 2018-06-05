# Read JPEG quality from bytes or os.File

Inpiration from [Estimating Quality](http://fotoforensics.com/tutorial-estq.php)


## usage:

````go

	file, err := os.Open("file.jpg")
	if err != nil {
		log.Fatal(err)
	}
	j, err := jpegquality.New(file) // or NewWithBytes([]byte)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("jpeg quality %d", j.Quality())
````
## logging:
There is some output to your current log. So, you can change output target with ```SetDefaultLoggerOutput``` before running code from above. E.g., disable logging:
````go

	jpegquality.SetDefaultLoggerOutput(ioutil.Discard)
	j, err := jpegquality.New(file) // or NewWithBytes([]byte)
	...
````
