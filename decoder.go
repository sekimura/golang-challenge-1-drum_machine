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
// TODO: implement
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

	scanp := 50
	tracks := make([]Track, 0)
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

		tracks = append(tracks, Track{
			Id:    id,
			Name:  str,
			Steps: steps,
		})
	}

	p := &Pattern{
		Version: version,
		Tempo:   120,
		Tracks:  tracks,
	}
	return p, nil
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
// TODO: implement
type Track struct {
	Id    int
	Name  string
	Steps []byte
}

type Pattern struct {
	Version string
	Tempo   int
	Tracks  []Track
}

func (p *Pattern) String() string {
	var b bytes.Buffer
	b.Write([]byte(fmt.Sprintf("Saved with HW Version: %s\n", p.Version)))
	b.Write([]byte(fmt.Sprintf("Tempo: %d\n", p.Tempo)))
	for _, t := range p.Tracks {
		b.Write([]byte(fmt.Sprintf("(%d) %s\t", t.Id, t.Name)))
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
