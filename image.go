package driverlima

import "github.com/kuttiproject/drivercore"

type Image struct {
	ImageK8sVersion string
	// imageChecksum   string
	ImageSourceURL  string
	ImageStatus     drivercore.ImageStatus
	ImageDeprecated bool
}

// K8sVersion returns the version of Kubernetes components in the image.
func (i *Image) K8sVersion() string {
	return i.ImageK8sVersion
}

// Status can be Notdownloaded, Downloaded or Unknown.
func (i *Image) Status() drivercore.ImageStatus {
	return i.ImageStatus
}

// Deprecated returns true if the image is no longer supported.
func (i *Image) Deprecated() bool {
	return i.ImageDeprecated
}

// Fetch downloads the image from the driver repository into the local cache.
// The lima driver does not download or cache the image; lima itself does
// that. So, Fetch silently continues if called.
func (i *Image) Fetch() error {
	i.ImageStatus = drivercore.ImageStatusDownloaded
	return nil
}

// FetchWithProgress downloads the image from the driver repository into the
// local cache, and reports progress via the supplied callback. The callback
// reports current and total in bytes.
// The lima driver does not download or cache the image; lima itself does
// that. So, FetchWithProgress silently continues if called.
func (i *Image) FetchWithProgress(progress func(current int64, total int64)) error {
	progress(100, 100)
	i.ImageStatus = drivercore.ImageStatusDownloaded
	return nil
}

// FromFile imports the image from the local filesystem into the local cache.
// The lima driver does not download or cache the image; lima itself does
// that. So, FromFile silently continues if called.
func (i *Image) FromFile(filepath string) error {
	i.ImageStatus = drivercore.ImageStatusDownloaded
	return nil
}

// PurgeLocal removes the image from the local cache.
// The lima driver does not download or cache the image; lima itself does
// that. So, PurgeLocal silently continues if called.
func (i *Image) PurgeLocal() error {
	i.ImageStatus = drivercore.ImageStatusNotDownloaded
	return nil
}
