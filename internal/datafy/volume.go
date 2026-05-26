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

func (c *Client) CreateVolume(ctx context.Context, req *CreateVolumeRequest) (*CreateVolumeResponse, error) {
	props := map[string]interface{}{
		"availabilityZone": req.AvailabilityZone,
		"diskSize":         req.DiskSize,
		"volumeIops":       req.VolumeIops,
		"volumeThroughput": req.VolumeThroughput,
		"encrypted":        req.Encrypted,
		"kmsKeyId":         req.KmsKeyId,
	}
	if len(req.Tags) > 0 {
		tags := make([]map[string]string, 0, len(req.Tags))
		for k, v := range req.Tags {
			tags = append(tags, map[string]string{"key": k, "value": v})
		}
		props["tags"] = tags
	}
	body := map[string]interface{}{"volumeProperties": props}
	resp, err := c.callAPI(ctx, http.MethodPost, "/api/v1/volumes/create-datafied-volume", body)
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
