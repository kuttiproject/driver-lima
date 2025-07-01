package driverlima

import (
	"encoding/json"
	"path"

	"github.com/kuttiproject/kuttilog"
	"github.com/kuttiproject/workspace"
)

const imagesConfigFile = "limaimages.json"

var (
	imagedata             = &imageconfigdata{}
	imageconfigmanager, _ = workspace.NewFileConfigManager(imagesConfigFile, imagedata)
)

type imageconfigdata struct {
	images map[string]*Image
}

func (icd *imageconfigdata) Serialize() ([]byte, error) {
	return json.Marshal(icd.images)
}

func (icd *imageconfigdata) Deserialize(data []byte) error {
	loaddata := make(map[string]*Image)
	err := json.Unmarshal(data, &loaddata)
	if err == nil {
		icd.images = loaddata
	}
	return err
}

func (icd *imageconfigdata) SetDefaults() {
	icd.images = defaultimages()
}

func defaultimages() map[string]*Image {
	return map[string]*Image{}
}

func limaConfigDir() (string, error) {
	return workspace.ConfigDir()
}

func fetchimagelist() error {
	// Download image list into temp file
	confdir, _ := limaConfigDir()
	tempfilename := "limaimagesnewlist.json"
	tempfilepath := path.Join(confdir, tempfilename)

	kuttilog.Printf(kuttilog.Debug, "confdir: %v\ntempfilepath: %v\n", confdir, tempfilepath)

	kuttilog.Println(kuttilog.Info, "Fetching image list...")
	kuttilog.Printf(kuttilog.Debug, "Fetching from %v into %v.", ImagesSourceURL, tempfilepath)
	err := workspace.DownloadFile(ImagesSourceURL, tempfilepath)
	kuttilog.Printf(kuttilog.Debug, "Error: %v", err)
	if err != nil {
		return err
	}
	defer workspace.RemoveFile(tempfilepath)

	configfilepath := path.Join(confdir, imagesConfigFile)

	err = workspace.CopyFile(tempfilepath, configfilepath, 32*1024, true)
	if err != nil {
		return err
	}

	// Load into object
	tempimagedata := &imageconfigdata{}
	tempconfigmanager, err := workspace.NewFileConfigManager(tempfilename, tempimagedata)
	if err != nil {
		return err
	}

	err = tempconfigmanager.Load()
	if err != nil {
		return err
	}

	// // Compare against current and update
	// for key, newimage := range tempimagedata.images {
	// 	oldimage := imagedata.images[key]
	// 	if oldimage != nil &&
	// 		// newimage.imageChecksum == oldimage.imageChecksum &&
	// 		newimage.imageSourceURL == oldimage.imageSourceURL &&
	// 		oldimage.imageStatus == drivercore.ImageStatusDownloaded {

	// 		newimage.imageStatus = drivercore.ImageStatusDownloaded
	// 	}
	// }

	// Make it current
	imagedata.images = tempimagedata.images

	// Save as local configuration
	imageconfigmanager.Save()

	return nil
}

func imagenamefromk8sversion(k8sversion string) string {
	return "kutti-k8s-" + k8sversion + ".qcow2"
}
