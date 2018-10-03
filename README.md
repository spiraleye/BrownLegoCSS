# BrownLegoCSS

This is a Go-based replacement for YUICompressor's CSS compressor.  This version seems to offer significant performance improvements over the original Java version, mostly due to it being compiled to native code, and also being able to compress/minify the contents of multiple files simultaneously through Go's awesome concurrency features.  However, there is probably a lot that can still be done in order to improve performance further, so please feel free to fork, improve, and submit pull requests!

## Files

* `./BrownLegoCSS.go` - Class/Lib/thing-that-does-the-hard-work.
* `./BrownLegoCSS_test.go` - Unit tests for the above.
* `./tests/*.css` - Input files for unit tests.
* `./tests/*.css.min` - Expected results for the above tests.
* `./csscompress/compress.go` - Go application that uses BrownLegoCSS lib to compress/minify CSS code.

## Name

Weird name, right?  Our inspiration for it comes from the GIF at http://thejoysofcode.com/post/37648725275/when-i-release-software-and-immediately-get-sued-by-a ;-)

## License

The copyrights embodied in the contents of this project's source files are licensed
by Spiraleye Studios under the 3-Clause BSD open source license,
as follows:

Copyright (c) 2012, Spiraleye Studios
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

* Neither the name of Spiraleye Studios nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL Spiraleye Studios BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

## Derivative Work

Most of the content in `BrownLegoCSS.go` is derived from the original Java code at [https://github.com/yui/yuicompressor/](https://github.com/yui/yuicompressor/).  For this reason, we have included the original license information below:

> YUI Compressor
[https://github.com/yui/yuicompressor/](http://developer.yahoo.com/yui/compressor/)
Author: Julien Lecomte -  [http://www.julienlecomte.net/](http://www.julienlecomte.net/)
Author: Isaac Schlueter - [http://foohack.com/](http://foohack.com/)
Author: Stoyan Stefanov - [http://phpied.com/](http://phpied.com/)
Copyright (c) 2011 Yahoo! Inc.  All rights reserved.
The copyrights embodied in the content of this file are licensed
by Yahoo! Inc. under the BSD (revised) open source license.
