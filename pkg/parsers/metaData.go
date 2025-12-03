package parsers

import (
	"bytes"
	"fmt"

	"github.com/dhowden/tag"
)

type FileInformation struct {
	Format  string
	Bitrate int
}

type Metadata struct {
	Name   string
	Artist string
	Album  string
	Year   int
	Format string
}

func (m Metadata) String() string {
	return fmt.Sprintf("Metadata{Name: %s, Artist: %s, Album: %s, Year: %d}", m.Name, m.Artist, m.Album, m.Year)
}

type AudioMetadata struct {
	SampleRate     int
	Channels       int
	DurationMs     int
	DecodedBitrate int
	Metadata       Metadata
}

func GetAudioMetadata(audioData *AudioData) (*AudioMetadata, error) {
	m, err := tag.ReadFrom(bytes.NewReader(audioData.RawData))

	if err != nil {
		return nil, err
	}

	name := m.Title()
	artist := m.Artist()
	album := m.Album()
	year := m.Year()
	format := m.Format()

	metadata := Metadata{
		Name:   name,
		Artist: artist,
		Album:  album,
		Year:   year,
		Format: string(format),
	}

	println("Name: ", metadata.String())

	decodedBitrate := int(audioData.Format.SampleRate) * audioData.Format.NumChannels * audioData.Format.Precision * 8

	return &AudioMetadata{
		SampleRate:     int(audioData.Format.SampleRate),
		Channels:       audioData.Format.NumChannels,
		DurationMs:     (int(audioData.Samples.Len()/int(audioData.Format.SampleRate)) * 1000),
		DecodedBitrate: decodedBitrate,
		Metadata:       metadata,
	}, nil
}
