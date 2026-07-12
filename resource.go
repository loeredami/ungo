package ungo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func PackResources(filePaths []string, outputPath string) error {
	buf := new(bytes.Buffer)

	// File count
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(filePaths))); err != nil {
		return err
	}

	for _, path := range filePaths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		name := []byte(path)

		if err := binary.Write(buf, binary.LittleEndian, uint16(len(name))); err != nil {
			return err
		}
		buf.Write(name)

		if err := binary.Write(buf, binary.LittleEndian, uint64(len(data))); err != nil {
			return err
		}
		buf.Write(data)
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

type Package struct {
	Files *SmallMap[string, []byte]
}

func LoadPackage(path string) (*Package, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	var count uint32
	if err := binary.Read(reader, binary.LittleEndian, &count); err != nil {
		return nil, err
	}

	pkg := &Package{
		Files: NewSmallMap[string, []byte](int(count)),
	}

	for i := uint32(0); i < count; i++ {
		var nameLen uint16
		if err := binary.Read(reader, binary.LittleEndian, &nameLen); err != nil {
			return nil, err
		}

		name := make([]byte, nameLen)
		if _, err := reader.Read(name); err != nil {
			return nil, err
		}

		var fileLen uint64
		if err := binary.Read(reader, binary.LittleEndian, &fileLen); err != nil {
			return nil, err
		}

		content := make([]byte, fileLen)
		if _, err := reader.Read(content); err != nil {
			return nil, err
		}

		pkg.Files.Set(string(name), content)
	}

	return pkg, nil
}

func (p *Package) Get(name string) ([]byte, bool) {
	return p.Files.Get(name)
}

type ResourceLoader struct {
	packages    []Lazy[*Package]
	loose_files SmallMap[string, Lazy[[]byte]]
}

func (r *ResourceLoader) Get(name string) Result[[]byte] {
	for _, pkg := range r.packages {
		if value, ok := pkg.Value().Get(name); ok {
			return VSuccess(value)
		}
	}

	if data, ok := r.loose_files.Get(name); ok {
		return VSuccess(data.Value())
	}

	if data, err := os.ReadFile(name); err == nil {
		r.loose_files.Set(name, NewLazy(func() []byte {
			d, _ := os.ReadFile(name)
			return d
		}))
		return VSuccess(data)
	}

	return VFail[[]byte](fmt.Errorf("resource not found %s", name))
}
