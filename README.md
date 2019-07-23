[![Build Status](https://travis-ci.org/bingoohuang/gou.svg?branch=master)](https://travis-ci.org/bingoohuang/gou)

# Go Utilities

## Strings

```go
fmt.Println(str.SingleLine("hello\nworld"))  // hello world
fmt.Println(str.FirstWord("hello world"))   // hello
fmt.Println(str.ParseMapString("k1=v1;k2=v2", ";", "="))   // map[k1:v1 k2:v2]
fmt.Println(str.MapOf("k1", "v1", "k2", "v2"))   // map[k1:v1 k2:v2]
fmt.Println(str.MapToString(map[string]string{"k1": "v1", "k2": "v2"})) // map[k1:v1 k2:v2]
fmt.Println(str.IndexOf("k1", "k0", "k1"))   // 1
fmt.Println(str.SplitTrim("k1,,k2", ","))   // [k1 k2]
fmt.Println(str.EmptyThen("", "default"))    // default
fmt.Println(str.ContainsIgnoreCase("ÑOÑO", "ñoño"))   // true
fmt.Println(str.HasPrefix("http://www.abc.com", "http://", "https://")) // true

a := ".tar.gz"
fmt.Println(str.AnyOf(a, ".tar", ".tar.gz")) // true
fmt.Println(str.NoneOf(a, ".xls", ".xlsx"))  // true

fmt.Println(str.ParseFloat32("1.1")) // 1.1
fmt.Println(str.ParseFloat64("1.1")) // 1.1
fmt.Println(str.ParseInt8("-11")) // -11
fmt.Println(str.ParseInt16("11")) // 11
fmt.Println(str.ParseInt32("11")) // 11
fmt.Println(str.ParseInt64("11")) // 11

fmt.Println(str.ParseUint("11")) // 11
fmt.Println(str.ParseUint8("11")) // 11
fmt.Println(str.ParseUint16("11")) // 11
fmt.Println(str.ParseUint32("11")) // 11
fmt.Println(str.ParseUint64("11")) // 11

fmt.Println(str.ParseInt("11")) // 11
fmt.Println(str.ParseInt8("11")) // 11
fmt.Println(str.ParseInt16("11")) // 11
fmt.Println(str.ParseInt32("11")) // 11
fmt.Println(str.ParseInt64("11")) // 11
```

## Codecs

1. Targz/Untargz

```go
fmt.Println(enc.Base64("不忘初心牢记使命!")) // 5LiN5b-Y5Yid5b-D54mi6K6w5L2_5ZG9IQ
fmt.Println(enc.Base64Decode("5LiN5b-Y5Yid5b-D54mi6K6w5L2_5ZG9IQ")) // 不忘初心牢记使命!

fmt.Println(enc.CBCEncrypt("16/24/32bytesxxx", "新时代中国特色社会主义!"))
fmt.Println(enc.CBCDecrypt("16/24/32bytesxxx", "HK5Ptmtt3V16mIBhJqNeQS_SbTn5kNmE4FSKoxx5t_I9fbIkf2GnjTF6T9KtuWuA8WZYWLMYZeAGsuHyycz9UA=="))
```

## Thanks

1. [A collection of useful utility functions](https://github.com/shomali11/util)
1. [github search go string utilities](https://github.com/search?l=Go&q=string+utilities&type=Repositories)
