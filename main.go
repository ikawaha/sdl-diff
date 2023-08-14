package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("specify two files")
	}
	src, err := loadSDL(args[0])
	if err != nil {
		return fmt.Errorf("src: %w", err)
	}
	dst, err := loadSDL(args[1])
	if err != nil {
		return fmt.Errorf("dst: %w", err)
	}

	if len(src) != len(dst) {
		fmt.Printf("len(src) != len(dst), %d <> %d\n", len(src), len(dst))
	} else {
		fmt.Printf("len(src) == len(dst), %d\n", len(src))
	}
	for k, v := range src {
		printDiff(v, dst[k])
	}
	return nil
}

func printDiff(lhs, rhs *Object) {
	if lhs == nil || rhs == nil {
		fmt.Printf("imcompatible object: %s, %s\n", lhs.Name(), rhs.Name())
		return
	}
	fmt.Printf("[%s] name: %s ==========\n", lhs.kind, lhs.name)
	if lhs.kind != rhs.kind {
		fmt.Println(">>kind:", lhs.kind)
		fmt.Println("<<kind:", rhs.kind)
	}
	if lhs.name != rhs.name {
		fmt.Println(">>name:", lhs.name)
		fmt.Println("<<name:", rhs.name)
	}
	for k := range lhs.item {
		if _, ok := rhs.item[k]; !ok {
			fmt.Println(">>item:", k)
		}
	}
	for k := range rhs.item {
		if _, ok := lhs.item[k]; !ok {
			fmt.Println("<<item:", k)
		}
	}
}

type Object struct {
	kind string
	name string
	item map[string]struct{}
}

func (o *Object) Name() string {
	if o == nil {
		return "<nil>"
	}
	return o.name
}

func (o Object) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "kind:%s, name:%s\n", o.kind, o.name)
	for k := range o.item {
		fmt.Fprintf(&b, "   %s\n", k)
	}
	return b.String()
}

func loadSDL(path string) (map[string]*Object, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(bytes.NewReader(b))
	ret := map[string]*Object{}
	var o *Object
	for n := 1; s.Scan(); n++ {
		l := strings.TrimSpace(s.Text())
		if l == "" {
			continue
		}
		if strings.HasSuffix(l, "{") {
			if o != nil {
				return nil, fmt.Errorf("line:%d, unexpected '}'", n)
			}
			o = &Object{
				kind: "unknown",
				item: map[string]struct{}{},
			}
			ts := strings.Split(l, " ")
			if len(ts) > 0 {
				o.kind = ts[0]
			}
			if len(ts) > 1 {
				o.name = ts[1]
			}
			continue
		}
		if l == "}" {
			if o == nil {
				return nil, fmt.Errorf("line:%d, unexpected '}'", n)
			}
			if o != nil {
				ret[o.name] = o
			}
			o = nil
			continue
		}
		if o != nil {
			if _, ok := o.item[l]; ok {
				return nil, fmt.Errorf("dup: kind=%s, name=%s, %s", o.kind, o.name, l)
			}
			o.item[l] = struct{}{}
		} else {
			fmt.Println("skip:", l)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}
