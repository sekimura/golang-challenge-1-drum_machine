package drum

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

var spliceHeader = "SPLICE"

func checkHeader(b []byte) error {
	if string(b[:len(spliceHeader)]) != spliceHeader {
		return fmt.Errorf("drum: not a splice file")
	}
	return nil
}

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := checkHeader(d); err != nil {
		return nil, fmt.Errorf("drum: not a splice file")
	}

	v := string(d[14:30])
	version := v[:strings.Index(string(v), "\x00")]

	// TODO: parse the Tempo part
	tempo := int(d[48] >> 1)

	scanp := 50
	var tracks []track
	for {
		if scanp >= len(d) {
			break
		}

		id := int(d[scanp])

		// get the size of track name first
		scanp += 4
		s := int(d[scanp])
		scanp++
		// and scan string with the size
		str := string(d[scanp : scanp+s])
		scanp += s

		steps := d[scanp : scanp+16]
		scanp += 16

		tracks = append(tracks, track{
			ID:    id,
			Name:  str,
			Steps: steps,
		})
	}

	p := &Pattern{
		Version: version,
		Tempo:   tempo,
		Tracks:  tracks,
	}
	return p, nil
}

type track struct {
	ID    int
	Name  string
	Steps []byte
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Version string
	Tempo   int
	Tracks  []track
}

func (p *Pattern) String() string {
	var b bytes.Buffer
	b.Write([]byte(fmt.Sprintf("Saved with HW Version: %s\n", p.Version)))
	b.Write([]byte(fmt.Sprintf("Tempo: %d\n", p.Tempo)))
	for _, t := range p.Tracks {
		b.Write([]byte(fmt.Sprintf("(%d) %s\t", t.ID, t.Name)))
		for i, s := range t.Steps {
			if i%4 == 0 {
				b.Write([]byte("|"))
			}
			if s == '\x00' {
				b.Write([]byte("-"))
			} else {
				b.Write([]byte("x"))
			}
		}
		b.Write([]byte("|\n"))
	}
	return b.String()
}
