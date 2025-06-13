package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/sangnt1552314/mangadex-tui/internal/api"
	"github.com/sangnt1552314/mangadex-tui/internal/models"
)

func GetColorStatus(status string) tcell.Color {
	switch status {
	case "ongoing":
		return tcell.ColorGreen
	case "completed":
		return tcell.ColorOrange
	case "hiatus":
		return tcell.ColorYellow
	case "cancelled":
		return tcell.ColorRed
	default:
		return tcell.ColorWhite
	}
}

func FormatTextYear(year int) string {
	if year == 0 {
		return "Unknown"
	}
	return strconv.Itoa(year)
}

func ShortenDescription(desc string, maxLength int) string {
	if desc == "" {
		return "No description available"
	}

	// Clean up newlines and extra spaces
	desc = strings.Join(strings.Fields(desc), " ")

	if len(desc) <= maxLength {
		return desc
	}

	// Find the last space before maxLength
	lastSpace := strings.LastIndex(desc[:maxLength], " ")
	if lastSpace == -1 {
		lastSpace = maxLength
	}

	// Truncate and add ellipsis
	return desc[:lastSpace] + "..."
}

func FormatTags(tags []models.Tag) string {
	if len(tags) == 0 {
		return "None"
	}

	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, fmt.Sprintf("[blue]%s[-]", tag.Attributes.Name["en"]))
	}
	return strings.Join(tagNames, " | ")
}

func FormatTextStatus(status string) string {
	switch status {
	case "ongoing":
		return "Ongoing"
	case "completed":
		return "Completed"
	case "hiatus":
		return "Hiatus"
	case "cancelled":
		return "Cancelled"
	default:
		return status
	}
}

func GetMangaImage(mangaID string, size int) image.Image {
	coverList, err := api.GetMangaCover(mangaID)
	if err != nil {
		log.Println("Error fetching cover for manga:", err)
		return nil
	}

	if len(coverList.Data) == 0 {
		return nil
	}

	coverURL := api.GetCoverURL(mangaID, coverList.Data[0].Attributes.FileName, size)
	resp, err := http.Get(coverURL)
	if err != nil {
		log.Println("Error fetching cover image:", err)
		return nil
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading cover image data:", err)
		return nil
	}

	contentType := http.DetectContentType(imgData)
	var img image.Image

	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(imgData))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(imgData))
	default:
		log.Printf("Unsupported image type: %s", contentType)
		return nil
	}

	if err != nil {
		log.Printf("Error decoding image: %v", err)
		return nil
	}

	return img
}
