package merger

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/restlesswhy/video-merger/internal/models"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Merge(files [2]string, mod models.Mod, outFileName string) error {
	first := ffmpeg.Input(files[0])
	second := ffmpeg.Input(files[1])

	switch mod {
	case models.PicInPic:
		w, _, err := getVideoSize(files[0])
		if err != nil {
			return err
		}

		size := float64(w) / 2.4
		s := strconv.Itoa(int(size))

		if err := ffmpeg.Filter([]*ffmpeg.Stream{first, second.Filter("scale", ffmpeg.Args{fmt.Sprintf("%s:-1", s)})}, "overlay", ffmpeg.Args{"main_w-overlay_w-10:main_h-overlay_h-10"}).Output(outFileName).OverWriteOutput().ErrorToStdOut().Run(); err != nil {
			return err
		}

	case models.SideBySide:
		if err := ffmpeg.Filter([]*ffmpeg.Stream{first, second}, "hstack", nil).Output(outFileName).OverWriteOutput().ErrorToStdOut().Run(); err != nil {
			return err
		}

	default:
		return errors.New("wrong video mod")
	}

	return nil
}

func getVideoSize(fileName string) (int, int, error) {
	data, err := ffmpeg.Probe(fileName)
	if err != nil {
		return 0, 0, errors.Wrap(err, "get video probe error")
	}

	type VideoInfo struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
			Width     int
			Height    int
		} `json:"streams"`
	}
	vInfo := &VideoInfo{}
	err = json.Unmarshal([]byte(data), vInfo)
	if err != nil {
		return 0, 0, errors.Wrap(err, "unmarshal probe error")
	}

	for _, s := range vInfo.Streams {
		if s.CodecType == "video" {
			return s.Width, s.Height, nil
		}
	}

	return 0, 0, errors.New("video codec is not found")
}
