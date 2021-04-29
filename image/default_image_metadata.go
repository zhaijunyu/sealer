package image

import (
	"context"
	"fmt"
	"gitlab.alibaba-inc.com/seadent/pkg/image/reference"
	imageutils "gitlab.alibaba-inc.com/seadent/pkg/image/utils"
	v1 "gitlab.alibaba-inc.com/seadent/pkg/types/api/v1"
	"sort"
)

type DefaultImageMetadataService struct {
	BaseImageManager
}

func (d DefaultImageMetadataService) Tag(imageName, tarImageName string) error {
	imageMetadataMap, err := imageutils.GetImageMetadataMap()
	if err != nil {
		return err
	}
	imageMetadata, ok := imageMetadataMap[imageName]
	if !ok {
		return fmt.Errorf("failed to found image %s", imageName)
	}
	imageMetadata.Name = tarImageName
	if err := imageutils.SetImageMetadata(imageMetadata); err != nil {
		return fmt.Errorf("failed to add tag %s, %s", tarImageName, err)
	}
	return nil
}

func (d DefaultImageMetadataService) List() ([]imageutils.ImageMetadata, error) {
	imageMetadataMap, err := imageutils.GetImageMetadataMap()
	if err != nil {
		return nil, err
	}
	var imageMetadataList []imageutils.ImageMetadata
	for _, imageMetadata := range imageMetadataMap {
		imageMetadataList = append(imageMetadataList, imageMetadata)
	}
	sort.Slice(imageMetadataList, func(i, j int) bool {
		return imageMetadataList[i].Name < imageMetadataList[j].Name
	})
	return imageMetadataList, nil
}

// invoke GetImage here
func (d DefaultImageMetadataService) GetImage(imageName string) (*v1.Image, error) {
	panic("not implemented")
}

func (d DefaultImageMetadataService) GetRemoteImage(imageName string) (v1.Image, error) {
	named, err := reference.ParseToNamed(imageName)
	if err != nil {
		return v1.Image{}, err
	}

	err = d.initRegistry(named.Domain())
	if err != nil {
		return v1.Image{}, err
	}

	manifest, err := d.registry.ManifestV2(context.Background(), named.Repo(), named.Tag())
	if err != nil {
		return v1.Image{}, err
	}

	return d.downloadImageManifestConfig(named, manifest.Config.Digest)
}