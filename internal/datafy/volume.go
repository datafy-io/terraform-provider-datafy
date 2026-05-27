package datafy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateVolumeRequest struct {
	AvailabilityZone string
	DiskSize         int64
	VolumeIops       int64
	VolumeThroughput int64
	Encrypted        bool
	KmsKeyId         string
	Tags             map[string]string
}

type CreateVolumeResponse struct {
	VolumeId        string   `json:"volumeId"`
	TargetVolumeIds []string `json:"targetVolumeIds"`
}

type volumeTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Volume struct {
	EbsId            string      `json:"volumeId"`
	AvailabilityZone string      `json:"availabilityZone"`
	DiskSize         uint64      `json:"diskSize"`
	Iops             uint32      `json:"iops"`
	Throughput       uint32      `json:"throughput"`
	Tags             []volumeTag `json:"tags"`
}

type getVolumesResponse struct {
	Volumes []Volume `json:"volumes"`
}

type createVolumeBody struct {
	VolumeProperties createVolumeProps `json:"volumeProperties"`
}

type createVolumeProps struct {
	AvailabilityZone string      `json:"availabilityZone"`
	DiskSize         int64       `json:"diskSize"`
	VolumeIops       int64       `json:"volumeIops"`
	VolumeThroughput int64       `json:"volumeThroughput"`
	Encrypted        bool        `json:"encrypted"`
	KmsKeyId         string      `json:"kmsKeyId,omitempty"`
	Tags             []volumeTag `json:"tags,omitempty"`
}

func (c *Client) CreateVolume(ctx context.Context, req *CreateVolumeRequest) (*CreateVolumeResponse, error) {
	props := createVolumeProps{
		AvailabilityZone: req.AvailabilityZone,
		DiskSize:         req.DiskSize,
		VolumeIops:       req.VolumeIops,
		VolumeThroughput: req.VolumeThroughput,
		Encrypted:        req.Encrypted,
		KmsKeyId:         req.KmsKeyId,
	}
	for k, v := range req.Tags {
		props.Tags = append(props.Tags, volumeTag{Key: k, Value: v})
	}
	resp, err := c.callAPI(ctx, http.MethodPost, "/api/v1/volumes/create-datafied-volume", createVolumeBody{VolumeProperties: props})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}
	var result CreateVolumeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVolume returns nil if the volume is not found.
func (c *Client) GetVolume(ctx context.Context, volumeId string) (*Volume, error) {
	resp, err := c.callAPI(ctx, http.MethodGet, fmt.Sprintf("/api/v1/volumes/%s", volumeId), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, toError(resp)
	}
	var result getVolumesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Volumes) == 0 {
		return nil, nil
	}
	return &result.Volumes[0], nil
}

func (c *Client) DeleteVolume(ctx context.Context, volumeId string) error {
	resp, err := c.callAPI(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/volumes/%s", volumeId), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return toError(resp)
	}
	return nil
}
