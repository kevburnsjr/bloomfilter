Bloom Filter
============

This bloom filter implementation is a Go port of https://github.com/jasondavies/bloomfilter.js

[![Go Report Card](https://goreportcard.com/badge/github.com/httpimp/bloomfilter?1)](https://goreportcard.com/report/github.com/httpimp/bloomfilter)
[![GoDoc](https://godoc.org/github.com/httpimp/bloomfilter?status.svg)](https://godoc.org/github.com/httpimp/bloomfilter)

The ability to build a bloom filter on the server in Go and evaluate that filter on the client in
Javascript can have immense value for comparing application state in distributed single page
applications with offline read/write capabilities and large data sets.

There are a lot of open source bloom filter implementations available on the internet.

For the most part, these implementations are not compatible. Every repo adds its own special
sauce to its hashing algorithms and hash derivation methods. After scouring the internet for several
hours searching for a multi language bloom filter implemenation, none appeared that fit the
requirements.

So, the most popular actively maintained javascript bloom filter on Github was selected and ported
to Go.

The reference implementation uses a non-standard FNV algorithm, but it is also less than 120 lines
of javascript. This project proves that the reference implementation can be easily ported by a
skilled developer to any desired language in less than a day.

### Go Example

```go
package main

import (
	"encoding/base64"
	"log"

	"github.com/httpimp/bloomfilter"
)

func main() {
	m, k := bloomfilter.EstimateParameters(10, 1e-6)
	bf := bloomfilter.New(m, k)
	bf.Add([]byte("foo"))
	bf.Add([]byte("bar"))
	encoded := base64.StdEncoding.EncodeToString(bf.ToBytes())
	fmt.Println(m)
	fmt.Println(k)
	fmt.Println(string(encoded))
}

```

> 288  
> 20  
> iCCACAiAACAACIgAAAIIAqCAAIgogCAIAIgACAIAigiAAIqA

### Javascript Example

Now we can take that same base64 encoded byte array and evaluate it with bloomfilter.js in the
browser

```js
var bits = sjcl.codec.base64.toBits("iCCACAiAACAACIgAAAIIAqCAAIgogCAIAIgACAIAigiAAIqA");
var bloom = new BloomFilter(bits, 20);
console.log(bloom.test("foo"));
console.log(bloom.test("bar"));
console.log(bloom.test("baz"));
bloom.add("baz");
console.log(sjcl.codec.base64.fromBits(bloom.buckets));
```

> true  
> true  
> false  
> iCCACAiAACAACIgAQAIIAqSIEKgowKEKAowIGBIgyomAAIqI

### Go Example Again

After deserializing the filter in javascript and altering it, we can send it back to the server
again to confirm that it now includes the additional element.

```go
package main

import (
	"encoding/base64"
	"log"

	"github.com/httpimp/bloomfilter"
)

func main() {
	decoded, err := base64.StdEncoding.DecodeString("iCCACAiAACAACIgAQAIIAqSIEKgowKEKAowIGBIgyomAAIqI")
	if err != nil {
		panic(err)
	}
	bf := bloomfilter.NewFromBytes(decoded, 21)
	log.Println(bf.Test([]byte("foo")))
	log.Println(bf.Test([]byte("bar")))
	log.Println(bf.Test([]byte("baz")))
	log.Println(bf.Test([]byte("bork")))
}

```

> true  
> true  
> true  
> false

### Standing on the shoulders of giants

Thanks to @jasondavies for creating the reference implementation, a functioning bloom filter
in < 120 lines of code. This port took about 6 hours.

Thanks to @willf for estimation.
