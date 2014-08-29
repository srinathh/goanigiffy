GoAniGiffy
==========

Animated GIFs have emerged as very convenient ways to post very short movie clips on the web. GoAniGiffy
is a small utility written in [Go language](www.golang.org) that converts a set of alphabetically sorted
images (eg. frames extracted from a video with mplayer or VLC) into an animated GIF with the ability to
**crop, scale, rotate & flip** the source images.

Requirements
------------
You need to have the Go language [installed](http://golang.org/doc/install) to build GoAniGiffy 

Installation
------------
You should be able to use go get to install GoAniGiffy. This will get the source and create a
binary built in your $GOPATH/bin folder
```
go get github.com
```

Usage
-----
GoAniGiffy performs image operations in the order of cropping, scaling, rotating & flipping before 
converting the images into an Animated GIF. Image manipulation is done using [Grigory Dryapak's imaging](www.github.com/disintegration/imaging)
package. We use the Lanczos filter in Resizing and the default Floyd-Steinberg dithering provided by
Go Language's [image/gif](http://golang.org/pkg/image/gif/) package to ensure video quality. 
Arbitrary angle rotations are not supported. 

The -delay parameter must be an integer specifying delay between frames in hundredths of a second. 
A value of 3 would give approximately 33 fps theoritically
```
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
```

Example
-------
Here is the command line that builds movie.gif from the images in the sample folder.
```
goanigiffy -src="sample/*.jpg" -dest="sample/movie.gif" -cropleft=100 -croptop=280 -cropwidth=550 -cropheight=351 -scale=0.5 -rotate=270 -verbose
```

Tips
----
Here is an example of how you can extract frames from a video clip using mplayer. We are extracting
JPEGs with quality of 80, starting from the 6th second & ending in the 8th second
```
mplayer -vo jpeg:quality=80 -nosound -ss 6 -endpos 8 vid.mp4
```
License
-------
GoAniGiffy code is licensed under the [Apache v2.0](https://github.com/srinathh/goanigiffy/blob/master/LICENSE) license.

All other media files & documentation in the repository is licened under <a rel="license" href="http://creativecommons.org/licenses/by-sa/4.0/">Creative Commons Attribution-ShareAlike 4.0 International License</a>.

