package jpegquality

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
)

var (
	ErrInvalidJPEG = errors.New("Invalid JPEG content")
)

type jpegReader struct {
	br      *bytes.Reader
	quality int
}

func NewWithBytes(buf []byte) (jr *jpegReader, err error) {
	jr = &jpegReader{br: bytes.NewReader(buf)}

	var (
		b1, b2 byte
	)
	b1, err = jr.br.ReadByte()
	if err != nil {
		return
	}
	b2, err = jr.br.ReadByte()
	if err != nil {
		return
	}
	if b1 != 0xff && b2 != 0xd8 {
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
			b1, b2        byte
			qualityAvg    = make([]float64, 3)
		)
		b1, err = this.br.ReadByte()
		if err != nil {
			log.Printf("read err %s", err)
			return
		}
		b2, err = this.br.ReadByte()
		if err != nil {
			log.Printf("read err %s", err)
			return
		}

		length = int(b1)*256 + int(b2) - 2
		if length < 0 {
			length = 0
		}
		// log.Printf("length %d from %x %x", length, b1, b2)

		if mark != 0xffdb { // not a quantization table
			b := make([]byte, length)
			var rt int
			rt, err = this.br.Read(b)
			if err != nil {
				log.Printf("read err %s", err)
				return
			} else {
				log.Printf("read %d bytes", rt)
			}
			continue
		}

		if length%65 != 0 {
			log.Printf("ERROR: Wrong size for quantization table -- this contains %d bytes (%d bytes short or %d bytes long)\n", length, 65-length%65, length%65)
		}

		log.Print("Quantization table")

		for length > 0 {
			var precision byte
			precision, err = this.br.ReadByte()
			if err != nil {
				log.Printf("read err %s", err)
				return
			}
			length--
			index = int(precision) & 0x0f
			precision = (precision & 0xf0) / 16
			log.Printf("  Precision=%d; Table index=%d (%s)\n", precision, index, getTableName(index))

			var (
				total, total_num int
			)

			for length > 0 && total_num < 64 {
				var b byte
				b, err = this.br.ReadByte()
				if err != nil {
					return
				}

				if total_num != 0 {
					total += int(b)
				} /* ignore first value */
				length--
				/* Show quantization table */
				if total_num%8 == 0 {
					fmt.Printf("    ")
				}
				fmt.Printf("%4d", b)
				if total_num%8 == 7 {
					fmt.Printf("\n")
				}
				total_num++
			}
			total_num--    /* we read 64 bytes, but only care about 63 values */
			if index < 3 { /* Only track the first 3 quantization tables */

				qualityAvg[index] = 100.0 - float64(total)/float64(total_num)
				fmt.Printf("  Estimated quality level = %5.2f%%\n", qualityAvg[index])
				if qualityAvg[index] <= 0 {
					fmt.Printf("  Quality too low; estimate may be incorrect.\n")
				}
				/* copy over the Q tables for initialization (in case Cr==Cb) */
				for i := index + 1; i < 3; i++ {
					qualityAvg[i] = qualityAvg[index]
				}
			}

			if index > 0 {
				// log.Printf("Averages(%d) %.2f %.2f %.2f", index, qualityAvg[0], qualityAvg[1], qualityAvg[2])
				var diff, qualityF float64
				diff = math.Abs(qualityAvg[0]-qualityAvg[1]) * 0.49
				diff += math.Abs(qualityAvg[0]-qualityAvg[2]) * 0.49
				qualityF = (qualityAvg[0]+qualityAvg[1]+qualityAvg[2])/3.0 + diff
				q = int(qualityF + 0.5)
				log.Printf("Average quality: %5.2f%% (%d%%)\n", qualityF, q)
				return
			}
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
		b1, b2 byte
		err    error
	)

ReadAgainB1:
	b1, err = this.br.ReadByte()
	if err != nil {
		return 0
	}
	if b1 != 0xff {
		b1, err = this.br.ReadByte()
		if err != nil {
			return 0
		}
	}

ReadAgainB2:
	b2, err = this.br.ReadByte()
	if err != nil {
		return 0
	}
	if b2 == 0xff {
		goto ReadAgainB2
	}
	if b2 == 0x00 {
		goto ReadAgainB1
	}

	// log.Printf("get marker %x %x", b1, b2)
	return int(b1)*256 + int(b2)
}

func (this *jpegReader) Quality() int {
	return this.quality
}
