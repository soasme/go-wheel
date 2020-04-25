package wheel

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"net/mail"
	"path/filepath"
	"regexp"
)

// The wheel filename is {distribution}-{version}(-{build tag})?-{python tag}-{abi tag}-{platform tag}.whl.
//
type WheelFile struct {
	// Distribution name, e.g. 'django', 'pyramid'.
	Distribution string
	// Distribution version, e.g. 1.0.
	Version string
	// Optional build number. Must start with a digit.
	BuildTag string
	// E.g. 'py27', 'py2', 'py3'.
	PythonTag string
	// E.g. 'cp33m', 'abi3', 'none'.
	ABITag string
	// E.g. 'linux_x86_64', 'any'.
	PlatformTag string
}

func ParseFilename(filename string) (f *WheelFile, err error) {
	filename = filepath.Base(filename)
	re := regexp.MustCompile(`^(?P<Distribution>[a-zA-Z0-9_-]+)-(?P<Version>[^-]+)-((?P<BuildTag>\d[^-]*)-)?(?P<PythonTag>[^-]+)-(?P<ABITag>[^-]+)-(?P<PlatformTag>[^-]+)\.whl$`)
	m := re.FindStringSubmatch(filename)
	if m == nil {
		return nil, fmt.Errorf("wheel: invalid file name convention: %s", filename)
	}
	f = &WheelFile{}
	for i, k := range re.SubexpNames() {
		switch k {
		case "Distribution":
			f.Distribution = m[i]
		case "Version":
			f.Version = m[i]
		case "BuildTag":
			f.BuildTag = m[i]
		case "PythonTag":
			f.PythonTag = m[i]
		case "ABITag":
			f.ABITag = m[i]
		case "PlatformTag":
			f.PlatformTag = m[i]
		}
	}
	return f, nil
}

type Metadata struct {
	Header mail.Header
	Body   string
}

func Open(filename string) (reader *zip.ReadCloser, err error) {
	return zip.OpenReader(filename)
}

func ReadMetadata(f *zip.File) (metadata *Metadata, err error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	msg, err := mail.ReadMessage(rc)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}

	metadata = &Metadata{Header: msg.Header, Body: string(body)}
	return metadata, nil
}

func (m *Metadata) FetchOne(key string) (string, bool) {
	v, ok := m.Header[key]
	if len(v) > 0 {
		return v[len(v)-1], true
	}
	return "", ok
}

func (m *Metadata) FetchAll(key string) ([]string, bool) {
	v, ok := m.Header[key]
	return v, ok
}

func FindFile(r *zip.ReadCloser, filename string) (f *zip.File, err error) {
	for _, f := range r.File {
		if f.Name == filename {
			return f, nil
		}
	}
	return nil, fmt.Errorf("wheel: file non-exists: %s", filename)
}

func (w *WheelFile) PathToDistInfo() string {
	return w.Distribution + "-" + w.Version + ".dist-info"
}

func (w *WheelFile) PathToWheel() string {
	return w.PathToDistInfo() + "/WHEEL"
}

func (w *WheelFile) PathToMetadata() string {
	return w.PathToDistInfo() + "/METADATA"
}
