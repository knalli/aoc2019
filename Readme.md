# Advent of Code 2019 Solutions
Welcome to my solutions of [Advent Of Code](http://adventofcode.com) 2019 (AOC 2019).

A huge thanks to @topaz and his team for providing this great service.  

## Disclaimer
These are my personal solutions of the Advent Of Code (AOC) 2019. The code is 
*not indented* to be perfect in any kind of area. This year, my personal 
competition was to learn more about Go. These snippets are here for everyone
learning more, too. 

If you think, there is a piece of improvement: Go to the code, 
fill a PR and we are all happy. Share the knowledge. 

## Structure
The AOC contains 24 days with at least one puzzle/question per day (mostly there are two parts).

* Base path is the root folder.
* Each day has has sub module named `day01`, `day02` until `day24`.
* The days `test001` and `test002` exist for testing only.
* The module `dayless` exists for explicit code sharing (common stuff).
* Depending on content, a day could import (exported) symbols of a (previous) day.

## Usage
First, ensure dependencies are installed.
 
For running the day `test001`
* CLI: just enter `go test001/main.go`
* IDE (like IntelliJ/Golang): just execute the `test001/main.go`.

## License / Copyright
Everything is free for all.

Licensed under MIT. Copyright Jan Philipp.