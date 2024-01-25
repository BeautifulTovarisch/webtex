# WebTeX

A tool for writing LaTeX in Markdown destined for the web.

## Background

I wrote WebTeX as a means of authoring mathematics and computer science notes
for my personal webpage. I found that other tools such as Sphinx and MathJax
did not quite meet all of my needs for diagramming, proofs etc.

After a [false start](https://github.com/BeautifulTovarisch/texxen), I decided
to write a proper program.

## Overview

WebTeX converts LaTeX embedded in Markdown files into SVGs and assembles the 
rendered output as an HTML document. Additionally, WebTex helps shorten the
feedback loop of authoring TeX by:

- Automatically re-rendering modified documents
- Serving rendered HTML from a static site

## Requirements

WebTex requires a TeX installation and [pdf2svg](https://github.com/dawbarton/pdf2svg). 
This project contains a [texlive profile](./texlive.profile) and [package list](./texlive.packages)
for reproducible installations.

At minimum, the following packages are recommended:

- pgf 
- pgfplots 
- amsmath 
- standalone 
- xcolor 
- bibtex

## Usage

## Contributing

### Getting Started

### Testing
