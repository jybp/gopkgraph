package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	help, trim         bool
	maxMods, maxStdlib int
	pkg                string
)

func main() {
	flag.StringVar(&pkg, "pkg", "", "Path to the package. Omit this flag to target the current directory.")
	flag.BoolVar(&trim, "trim", true, "Trim module prefix.")
	flag.IntVar(&maxMods, "mods", 0, "Max depth for packages from other modules.")
	flag.IntVar(&maxStdlib, "stdlib", 0, "Max depth for packages from the stdlib.")
	flag.BoolVar(&help, "h", false, "Help.")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

var (
	// stdpkgs is a map of all stdlib packages.
	stdpkgs = make(map[string]struct{})

	// deps is a map of all deps that have been seen.
	deps = make(map[string]struct{})

	stdout io.Writer = os.Stdout
)

func run() error {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName}, "std")
	if err != nil {
		return fmt.Errorf("could not load stdlib packages: %v", err)
	}
	stdpkgs = make(map[string]struct{})
	for _, pkg := range pkgs {
		stdpkgs[pkg.PkgPath] = struct{}{}
	}
	var cfg = &packages.Config{
		Mode: packages.NeedImports | packages.NeedName | packages.NeedModule,
		Dir:  pkg,
	}
	pkgs, err = packages.Load(cfg, "")
	if err != nil {
		return fmt.Errorf("could not load package '%s': %v", pkg, err)
	}
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].PkgPath < pkgs[j].PkgPath })
	for _, pkg := range pkgs {
		if err := imports(pkg, pkg.Module.Path, 0, 0); err != nil {
			return err
		}
	}
	return nil
}

func imports(pkg *packages.Package, module string, stdlibDepth, modsDepth int) error {
	imps := []*packages.Package{}
	for _, imp := range pkg.Imports {
		imps = append(imps, imp)
	}
	sort.Slice(imps, func(i, j int) bool { return imps[i].PkgPath < imps[j].PkgPath })
	for _, imp := range imps {
		_, isStdlib := stdpkgs[imp.PkgPath]
		if isStdlib && stdlibDepth >= maxStdlib {
			continue
		}
		if isStdlib && pkg.Module != nil && pkg.Module.Path != module {
			continue
		}

		isOtherModule := imp.Module != nil && imp.Module.Path != module
		if isOtherModule && modsDepth >= maxMods {
			continue
		}

		if _, ok := deps[fmt.Sprintf("%s->%s", pkg.PkgPath, imp.PkgPath)]; ok {
			continue
		}
		deps[fmt.Sprintf("%s->%s", pkg.PkgPath, imp.PkgPath)] = struct{}{}

		if trim {
			pkg.PkgPath = strings.TrimPrefix(pkg.PkgPath, module+"/")
			imp.PkgPath = strings.TrimPrefix(imp.PkgPath, module+"/")
		}
		fmt.Fprintf(stdout, "%s %s\n", pkg.PkgPath, imp.PkgPath)

		var err error
		if isStdlib {
			err = imports(imp, module, stdlibDepth+1, modsDepth)
		} else if isOtherModule {
			err = imports(imp, module, stdlibDepth, modsDepth+1)
		} else {
			err = imports(imp, module, stdlibDepth, modsDepth)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
