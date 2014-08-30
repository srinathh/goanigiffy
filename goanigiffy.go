/*
   Copyright 2014 Hariharan Srinath

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

/*
GoAniGiffy is a utility for converting a set of alphabetically sorted images such as video frames
grabbed from VLC or MPlayer into an animated GIF with options to Crop, Resize, Rotate & Flip the
images prior to creating the GIF

GoAniGiffy performs image operations in the order of cropping, scaling, rotating & flipping
before converting the images into an Animated GIF. Image manipulation is done using
Grigory Dryapak's imaging package. We use the Lanczos filter in Resizing and the default
Floyd-Steinberg dithering used by Go Language's image/gif package to ensure video quality.
Arbitrary angle rotations are not supported.

The -delay parameter must be an integer specifying delay between frames in hundredths of
a second. A value of 3 would give approximately 33 fps theoritically

Usage of goanigiffy:
  -cropheight=-1: height of cropped image, -1 specified full height
  -cropleft=0: left co-ordinate for crop to start
  -croptop=0: top co-ordinate for crop to start
  -cropwidth=-1: width of cropped image, -1 specifies full width
  -delay=3: delay time between frame in hundredths of a second
  -dest="movie.gif": a destination filename for the animated gif
  -flip="none": valid falues are none, horizontal, vertical
  -rotate=0: valid values are 0, 90, 180, 270
  -scale=1: scaling factor to apply if any
  -src="*.jpg": a glob pattern for source images. defaults to *.jpg
  -verbose=false: show in-process messages

Sources: https://github.com/srinathh/goanigiffy
*/
package main

import (
	"bytes"
	"flag"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/disintegration/imaging"
)

func CropImage(cropleft, croptop, cropwidth, cropheight int, img image.Image, verbose bool) image.Image {
	//Crop operation. Ignore if there is no crop operation specified
	if !(cropwidth == -1 && cropheight == -1 && cropleft == 0 && croptop == 0) {
		if cropwidth == -1 {
			cropwidth = img.Bounds().Dx()
		}
		if cropheight == -1 {
			cropheight = img.Bounds().Dy()
		}
		if verbose {
			log.Printf("Cropping original image at (%d,%d)->(%d,%d)", cropleft, croptop, cropleft+cropwidth-1, croptop+cropheight-1)
		}
		img = imaging.Crop(img, image.Rect(cropleft, croptop, cropleft+cropwidth-1, croptop+cropheight-1))
	}
	return img
}

func ScaleImage(scale float64, img image.Image, verbose bool) image.Image {
	//Scale operation. Ignore if scale is 1.0
	if scale != 1.0 {
		newwidth := int(float64(img.Bounds().Dx()) * scale)
		newheight := int(float64(img.Bounds().Dy()) * scale)

		if verbose {
			log.Printf("Scaling image from (%d, %d) -> (%d, %d)", img.Bounds().Dx(), img.Bounds().Dy(), newwidth, newheight)
		}
		img = imaging.Resize(img, newwidth, newheight, imaging.Lanczos)
	}
	return img

}

func RotateImage(rotate int, img image.Image, verbose bool) image.Image {
	//Rotate operation. Ignore if rotate is 0
	if rotate != 0 && verbose {
		log.Printf("Rotating by %d", rotate)
	}
	switch rotate {
	case 90:
		img = imaging.Rotate90(img)
	case 180:
		img = imaging.Rotate180(img)
	case 270:
		img = imaging.Rotate270(img)
	}
	return img
}

//FlipImage takes a string
func FlipImage(flip string, img image.Image, verbose bool) image.Image {
	//Flip operation
	if flip != "none" && verbose {
		log.Printf("Flipping %s", flip)
	}

	switch flip {
	case "horizontal":
		img = imaging.FlipH(img)
	case "vertical":
		img = imaging.FlipV(img)
	}
	return img
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	srcglob := flag.String("src", "*.jpg", "a glob pattern for source images. defaults to *.jpg")
	destname := flag.String("dest", "movie.gif", "a destination filename for the animated gif")
	cropleft := flag.Int("cropleft", 0, "left co-ordinate for crop to start")
	croptop := flag.Int("croptop", 0, "top co-ordinate for crop to start")
	cropwidth := flag.Int("cropwidth", -1, "width of cropped image, -1 specifies full width")
	cropheight := flag.Int("cropheight", -1, "height of cropped image, -1 specified full height")
	delay := flag.Int("delay", 3, "delay time between frame in hundredths of a second")
	verbose := flag.Bool("verbose", false, "show in-process messages")
	scale := flag.Float64("scale", 1.0, "scaling factor to apply if any")
	rotate := flag.Int("rotate", 0, "valid values are 0, 90, 180, 270")
	flip := flag.String("flip", "none", "valid falues are none, horizontal, vertical")

	flag.Parse()

	if !(*rotate == 0 || *rotate == 90 || *rotate == 180 || *rotate == 270) {
		log.Printf("rotate flag must be one of 0, 90, 180 or 270")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !(*flip == "none" || *flip == "horizontal" || *flip == "vertical") {
		log.Printf("flip flag must be one of none, horizontal or vertical")
		flag.PrintDefaults()
		os.Exit(1)
	}

	srcfilenames, err := filepath.Glob(*srcglob)
	if err != nil {
		log.Fatalf("Error in globbing source file pattern %s : %s", *srcglob, err)
	}

	if len(srcfilenames) == 0 {
		log.Fatalf("No source images found via pattern %s", *srcglob)
	}

	if *verbose {
		log.Printf("Found %d images to parse", len(srcfilenames))
	}

	sort.Strings(srcfilenames)

	var frames []*image.Paletted

	for ctr, filename := range srcfilenames {
		img, err := imaging.Open(filename)
		if err != nil {
			log.Printf("Skipping file %s due to error reading it :%s", filename, err)
			continue
		}

		if *verbose {
			log.Printf("Parsing image %d of %d : %s", ctr, len(srcfilenames), filename)
		}

		img = CropImage(*cropleft, *croptop, *cropwidth, *cropheight, img, *verbose)
		img = ScaleImage(*scale, img, *verbose)
		img = RotateImage(*rotate, img, *verbose)
		img = FlipImage(*flip, img, *verbose)

		buf := bytes.Buffer{}
		if err := gif.Encode(&buf, img, nil); err != nil {
			log.Printf("Skipping file %s due to error in gif encoding:%s", filename, err)
			continue
		}

		tmpimg, err := gif.Decode(&buf)
		if err != nil {
			log.Printf("Skipping file %s due to weird error reading the temporary gif :%s", filename, err)
			continue
		}
		frames = append(frames, tmpimg.(*image.Paletted))
	}

	if *verbose {
		log.Printf("Parsed all images.. now attemting to create animated GIF %s", *destname)
	}

	delays := make([]int, len(frames))
	for j, _ := range delays {
		delays[j] = *delay
	}

	opfile, err := os.Create(*destname)
	if err != nil {
		log.Fatalf("Error creating the destination file %s : %s", *destname, err)
	}

	if err := gif.EncodeAll(opfile, &gif.GIF{frames, delays, 0}); err != nil {
		log.Printf("Error encoding output into animated gif :%s", err)
	}
	opfile.Close()
}
