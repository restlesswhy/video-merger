package v1

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/restlesswhy/video-merger/internal/merger"
	"github.com/restlesswhy/video-merger/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type App interface {
}

type handler struct {
	app App
}

func New(app App) *handler {
	return &handler{app: app}
}

func (h *handler) uploadVideo(c *fiber.Ctx) error {
	id := string(c.Context().FormValue("id"))
	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "id parameter is empty",
		})
	}

	userID := string(c.Context().FormValue("user_id"))
	if len(userID) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "user_id parameter is empty",
		})
	}

	videoID := string(c.Context().FormValue("video_id"))
	if len(userID) == 0 || videoID != "1" && videoID != "2" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "video_id parameter must be 1 or 2",
		})
	}

	var b strings.Builder
	b.WriteString("tmp/in_")
	b.WriteString(id)
	b.WriteString(userID)
	b.WriteString(videoID)
	b.WriteString(".mp4")

	file, err := os.Create(b.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"err":    err.Error(),
		})
	}
	defer file.Close()

	r := bytes.NewReader(c.Body())
	_, err = io.Copy(file, r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"err":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

func (h *handler) getMergedVideo(c *fiber.Ctx) error {
	id := string(c.Context().FormValue("id"))
	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "id parameter is empty",
		})
	}

	userID := string(c.Context().FormValue("user_id"))
	if len(userID) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "user_id parameter is empty",
		})
	}

	modStr := string(c.Context().FormValue("mod"))
	if len(modStr) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "mod parameter is empty",
		})
	}

	mod := models.Mod(modStr)
	if err := mod.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"err":    "wrong mod parameter",
		})
	}

	var files [2]string

	var b strings.Builder
	for i := 0; i < 2; i++ {
		b.WriteString("tmp/in_")
		b.WriteString(id)
		b.WriteString(userID)
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(".mp4")

		if _, err := os.Stat(b.String()); errors.Is(err, os.ErrNotExist) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"err":    "not enough video to concatenate",
			})
		}

		files[i] = b.String()

		b.Reset()
	}

	outFileName := fmt.Sprintf("tmp/%s%x.mp4", mod, sha256.Sum256([]byte(files[0])))

	if _, err := os.Stat(outFileName); errors.Is(err, os.ErrNotExist) {
		err := merger.Merge(files, mod, outFileName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"err":    err.Error(),
			})
		}
	}

	return c.SendFile(outFileName)
}
