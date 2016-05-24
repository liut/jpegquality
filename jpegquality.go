package jpegquality

import (
	"bytes"
	"errors"
	"io"
	"log"
	"math"
)

var (
	ErrInvalidJPEG = errors.New("Invalid JPEG content")
	ErrWrongTable  = errors.New("ERROR: Wrong size for quantization table")
)

type jpegReader struct {
	rs      io.ReadSeeker
	quality int
}

func NewWithBytes(buf []byte) (jr *jpegReader, err error) {
	return New(bytes.NewReader(buf))
}

func New(rs io.ReadSeeker) (jr *jpegReader, err error) {
	jr = &jpegReader{rs: rs}

	var (
		sign = make([]byte, 2)
	)
	_, err = jr.rs.Read(sign)
	if err != nil {
		return
	}
	if sign[0] != 0xff && sign[1] != 0xd8 {
		err = ErrInvalidJPEG
	}

	var q int
	q, err = jr.readQuality()
	if err == nil {
		jr.quality = q
	}
	return
}

func (this *jpegReader) readQuality() (q int, err error) {
	for {
		mark := this.readMarker()
		if mark == 0 {
			err = ErrInvalidJPEG
			return
		}
		var (
			length, index int
			sign          = make([]byte, 2)
			qualityAvg    = make([]float64, 3)
		)
		_, err = this.rs.Read(sign)
		if err != nil {
			log.Printf("read err %s", err)
			return
		}

		length = int(sign[0])*256 + int(sign[1]) - 2
		if length < 0 {
			length = 0
		}

		if mark != 0xffdb { // not a quantization table
			_, err = this.rs.Seek(int64(length), 1)
			if err != nil {
				log.Printf("seek err %s", err)
				return
			}
			continue
		}

		if length%65 != 0 {
			log.Printf("ERROR: Wrong size for quantization table -- this contains %d bytes (%d bytes short or %d bytes long)\n", length, 65-length%65, length%65)
			err = ErrWrongTable
			return
		}

		log.Printf("length %d", length)
		log.Print("Quantization table")

		var tabuf = make([]byte, length)
		_, err = this.rs.Read(tabuf)
		if err != nil {
			log.Printf("read err %s", err)
			return
		}
		for j := 0; j < int(float64(length)/float64(65)); j++ {
			buf := tabuf[j*65 : (j+1)*65]
			index = int(buf[0] & 0x0f)
			precision := (buf[0] & 0xf0) / 16
			log.Printf("  Precision=%d; Table index=%d (%s)\n", precision, index, getTableName(index))

			var total int
			for i, b := range buf {
				if i > 1 { // ignore first value and index
					total += int(b)
				}
			}
			log.Printf("total %d", total)
			qualityAvg[index] = 100.0 - float64(total)/63.0
		}

		if index < 3 {
			for i := index + 1; i < 3; i++ {
				qualityAvg[i] = qualityAvg[index]
			}
		}

		if index > 0 {
			log.Printf("Averages(%d) %.2f %.2f %.2f", index, qualityAvg[0], qualityAvg[1], qualityAvg[2])
			var diff, qualityF float64
			diff = math.Abs(qualityAvg[0]-qualityAvg[1]) * 0.49
			diff += math.Abs(qualityAvg[0]-qualityAvg[2]) * 0.49
			qualityF = (qualityAvg[0]+qualityAvg[1]+qualityAvg[2])/3.0 + diff
			q = int(qualityF + 0.5)
			log.Printf("Average quality: %5.2f%% (%d%%)\n", qualityF, q)
			return
		}
	}
	return
}

func getTableName(index int) string {
	if index > 0 {
		return "chrominance"
	}
	return "luminance"
}

func (this *jpegReader) readMarker() int {
	var (
		mark = make([]byte, 2)
		err  error
	)

ReadAgain:
	_, err = this.rs.Read(mark)
	if err != nil {
		return 0
	}
	if mark[0] != 0xff || mark[1] == 0xff || mark[1] == 0x00 {
		goto ReadAgain
	}

	// log.Printf("get marker %x", mark)
	return int(mark[0])*256 + int(mark[1])
}

func (this *jpegReader) Quality() int {
	return this.quality
}
