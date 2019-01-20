package filesystem

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	humanize "github.com/dustin/go-humanize"
)

type DiskUsageReporter struct {
}

type DiskUsageRequest struct {
	Root         string
	MaxRecursion int
	Verbose      bool
}
type DiskUsageResponse struct {
	Root *Dir
}

func (r *DiskUsageReporter) Do(ctx context.Context, req *DiskUsageRequest) (*DiskUsageResponse, error) {
	root, err := makeDirTree(ctx, req.Root, req.MaxRecursion, 0, req.Verbose)
	if err != nil {
		return nil, err
	}
	return &DiskUsageResponse{
		Root: root,
	}, nil
}

type Dir struct {
	Path  string
	Files []os.FileInfo
	Dirs  []*Dir
}

func (d *Dir) Size() uint64 {
	var s uint64
	for _, info := range d.Files {
		s += uint64(info.Size())
	}
	for _, dir := range d.Dirs {
		s += uint64(dir.Size())
	}
	return s
}

type DiskUsageEntry struct {
	Path  string
	Level int
	Size  uint64
}

func (e *DiskUsageEntry) String() string {
	return fmt.Sprintf("%s %s", e.Path, humanize.Bytes(e.Size))
}

type DiskUsageEntries struct {
	entries []*DiskUsageEntry
	root    string

	level      int
	order      string
	humanize   bool
	ibytes     bool
	printEmpty bool
	relative   bool
}

func (es *DiskUsageEntries) Level(level int) *DiskUsageEntries {
	es.level = level
	return es
}

func (es *DiskUsageEntries) Order(order string) *DiskUsageEntries {
	es.order = order
	return es
}

func (es *DiskUsageEntries) Humanize(humanize bool) *DiskUsageEntries {
	es.humanize = humanize
	return es
}

func (es *DiskUsageEntries) IBytes(ibytes bool) *DiskUsageEntries {
	es.ibytes = ibytes
	return es
}

func (es *DiskUsageEntries) Relative(relative bool) *DiskUsageEntries {
	es.relative = relative
	return es
}

func (es *DiskUsageEntries) Find() []*DiskUsageEntry {
	var found []*DiskUsageEntry
	predicator := es.predicator()
	for _, e := range es.entries {
		if predicator(e) {
			found = append(found, e)
		}
	}
	switch es.order {
	case "", "size":
		sort.Slice(found, func(i, j int) bool {
			return found[i].Size > found[j].Size
		})
	}
	return found
}

func (es *DiskUsageEntries) Print(w io.Writer) {
	founds := es.Find()
	var longestPath string
	for _, found := range founds {
		if len(found.Path) > len(longestPath) {
			longestPath = found.Path
		}
	}
	for _, found := range founds {
		if !es.printEmpty && found.Size == 0 {
			continue
		}
		var size string
		if es.ibytes {
			size = humanize.IBytes(found.Size)
		} else if es.humanize {
			size = humanize.Bytes(found.Size)
		} else {
			size = fmt.Sprintf("%d", found.Size)
		}

		path := found.Path
		if es.relative {
			if rel, err := filepath.Rel(es.root, path); err == nil {
				path = rel
			}
		}
		fmt.Fprintf(w, "%-*s\t%s\n", len(longestPath), path, size)
	}
}

func (es *DiskUsageEntries) predicator() func(*DiskUsageEntry) bool {
	return func(e *DiskUsageEntry) bool {
		if es.level != -1 && es.level < e.Level {
			return false
		}
		return true
	}
}

func (d *Dir) Entries() *DiskUsageEntries {
	return &DiskUsageEntries{
		entries: d.entries(-1),
		root:    d.Path,
	}
}

func (d *Dir) entries(level int) []*DiskUsageEntry {
	level += 1
	var entries []*DiskUsageEntry

	entries = append(entries, d.entry(level))
	for _, dir := range d.Dirs {
		entries = append(entries, dir.entries(level)...)
	}
	return entries
}

func (d *Dir) entry(level int) *DiskUsageEntry {
	return &DiskUsageEntry{
		Path:  d.Path,
		Level: level,
		Size:  d.Size(),
	}
}

func (d *Dir) Dump(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s %s\n", d.Path, humanize.Bytes(d.Size())); err != nil {
		return err
	}
	for _, dir := range d.Dirs {
		if err := dir.Dump(w); err != nil {
			return err
		}
	}
	return nil
}

var errSkip = errors.New("skip")

func makeDirTree(ctx context.Context, root string, maxLevel, currentLevel int, verbose bool) (*Dir, error) {
	log := func(msg string) {
		if verbose {
			fmt.Println(msg)
		}
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if maxLevel != -1 && currentLevel > maxLevel {
		log(fmt.Sprintf("skip %s current: %d max: %d", root, currentLevel, maxLevel))
		return nil, errSkip
	}

	log(fmt.Sprintf("search %s", root))
	infos, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	d := &Dir{Path: root}

	for _, info := range infos {
		if info.IsDir() {
			dir, err := makeDirTree(ctx, filepath.Join(root, info.Name()), maxLevel, currentLevel+1, verbose)
			if err == errSkip {
				continue
			} else if err != nil {
				return nil, err
			}
			d.Dirs = append(d.Dirs, dir)
		}

		d.Files = append(d.Files, info)
	}

	return d, nil
}
