package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var rootDir = "."
var count = 1000

func main() {

	start := time.Now()
	log.Printf("Start Time: %+v\n", start)

	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go imageWork(&wg, i)
	}
	wg.Wait()

	d := time.Since(start)
	log.Printf("Time to complete: %+v seconds\n", d.Seconds())
}

func imageWork(wg *sync.WaitGroup, index int) {
	img := makeImage(index)
	saveImage(img, index)
	wg.Done()
}

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func makeImage(index int) image.Image {
	r := image.Rect(0, 0, 1024, 960)
	m := image.NewRGBA(r)
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	log.Printf("Making new image in memory: %d\n", index)
	return m.SubImage(r)
}

func saveImage(img image.Image, index int) error {
	fpath := fmt.Sprintf(rootDir+string(os.PathSeparator)+"%04d.jpg", index)
	f, err := os.Create(fpath)
	if err != nil {
		log.Fatalf("Could not create output file %s: %+v\n", fpath, err)
	}
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	if err == nil {
		log.Printf("Saved new image: %s\n", fpath)
	} else {
		log.Printf("Error! saving new image: %s\n", fpath)
	}
	return err
}

func removeAllImages() {
	matches, _ := filepath.Glob(rootDir + string(os.PathSeparator) + "*.jpg")
	for _, v := range matches {
		os.Remove(v)
	}
}
