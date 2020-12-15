# LINE TechTrain - The example of TODO App

:earth_americas: [日本語版](README.ja.md)

## Overview

This repository is a reference material for the TechTrain mission from LINE. This repository consists of these files.

- /example
	- /client
		- The implementation example of FrontEnd. Written in HTML/CSS/JS
	- /server
		- The implementation example of Application Server. Written in Go.
- /wireframe
	- /toppage.png
		- The wireframe for "Top page"
	- /detail.png
		- The wireframe for "Detail page"
- /architecture
	- system.png
		- the brief system overview

The execution environment for these is not provided.

## The purpose of these references

It's important to grasp an overview before getting into development phase.
The contents in this repository are just referencies to understand the overview of the TechTrain mission. Please regard them as just an example.
We don't have the intention that you implement exactly as the referencies said.

## The example of each page's requirement

- Top page
	- Show a list TODO items
	- Each item is a link to "Detail page" and has a "finish button" in the left of it
	- A list of finished items. The list is collapsable and closed as default. Clicking it, show the items
- Detail page
	- Show a form
		- an input form for TODO’s title
		- an input form for TODO’s description
		- a register button
		- a delete button

## See also

- [line/line-liff-v2-starter](https://github.com/line/line-liff-v2-starter)
	- The repository showing the example of how to develop LIFF app

## License

You can use the contents in this repository only when you work on our TechTrain mission.
Reproduction is prohibited.