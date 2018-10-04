package mongodbatlas

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterService_List(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/atlas/v1.0/groups/123/clusters", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		fmt.Fprintf(w, `{"links":[],"results":[{"name":"test","mongoDBMajorVersion":"3.4"}],"totalCount":1}`)
	})

	client := NewClient(httpClient)
	clusters, _, err := client.Clusters.List("123")
	expected := []Cluster{Cluster{Name: "test", MongoDBMajorVersion: "3.4"}}
	assert.Nil(t, err)
	assert.Equal(t, expected, clusters)
}

func TestClusterService_Get(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/atlas/v1.0/groups/123/clusters/test", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "GET", r)
		fmt.Fprintf(w, `{"name":"test","mongoDBMajorVersion":"3.4"}`)
	})

	client := NewClient(httpClient)
	cluster, _, err := client.Clusters.Get("123", "test")
	expected := &Cluster{Name: "test", MongoDBMajorVersion: "3.4"}
	assert.Nil(t, err)
	assert.Equal(t, expected, cluster)
}

func TestClusterService_Create(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/atlas/v1.0/groups/123/clusters", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "POST", r)
		w.Header().Set("Content-Type", "application/json")
		expectedBody := map[string]interface{}{
			"name":                  "test",
			"mongoDBMajorVersion":   "3.4",
			"replicationFactor":     float64(3),
			"backupEnabled":         false,
			"providerBackupEnabled": false,
			"paused":                false,
			"diskSizeGB":            10.5,
			"autoScaling": map[string]interface{}{
				"diskGBEnabled": false,
			},
			"providerSettings": map[string]interface{}{
				"providerName":     "AWS",
				"regionName":       "US_EAST_1",
				"instanceSizeName": "M0",
			},
			"replicationSpec": map[string]interface{}{
				"US_EAST_1": map[string]interface{}{
					"priority":       float64(7),
					"electableNodes": float64(2),
					"readOnlyNodes":  float64(1),
				},
			},
		}
		assertReqJSON(t, expectedBody, r)
		fmt.Fprintf(w, `{
			"name":"test",
			"mongoDBMajorVersion":"3.4",
			"replicationFactor":3,
			"backupEnabled":false,
			"providerBackupEnabled":false,
			"diskSizeGB":10,
			"paused":false,
			"autoScaling":{
				"diskGBEnabled":false
			},
			"providerSettings":{
				"providerName":"AWS",
				"regionName":"US_EAST_1",
				"instanceSizeName":"M0"
			},
			"replicationSpec":{
				"US_EAST_1":{
					"priority":7,
					"electableNodes":2,
					"readOnlyNodes":1
				}
			}
		}`)
	})

	client := NewClient(httpClient)
	providerSettings := ProviderSettings{ProviderName: "AWS", RegionName: "US_EAST_1", InstanceSizeName: "M0"}
	replicationSpec := map[string]ReplicationSpec{
		"US_EAST_1": ReplicationSpec{Priority: 7, ElectableNodes: 2, ReadOnlyNodes: 1},
	}
	params := &Cluster{
		Name:                "test",
		MongoDBMajorVersion: "3.4",
		ReplicationFactor:   3,
		BackupEnabled:       false,
		DiskSizeGB:          10.5,
		ProviderSettings:    providerSettings,
		ReplicationSpec:     replicationSpec,
	}
	cluster, _, err := client.Clusters.Create("123", params)
	expected := &Cluster{
		Name:                "test",
		MongoDBMajorVersion: "3.4",
		ReplicationFactor:   3,
		BackupEnabled:       false,
		Paused:              false,
		DiskSizeGB:          10,
		ProviderSettings:    providerSettings,
		ReplicationSpec:     replicationSpec,
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, cluster)
}

func TestClusterService_Update(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/atlas/v1.0/groups/123/clusters/test", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "PATCH", r)
		w.Header().Set("Content-Type", "application/json")
		expectedBody := map[string]interface{}{
			"diskSizeGB":            float64(5),
			"backupEnabled":         false,
			"providerBackupEnabled": false,
			"paused":                false,
			"autoScaling": map[string]interface{}{
				"diskGBEnabled": false,
			},
			"providerSettings": map[string]interface{}{},
		}
		assertReqJSON(t, expectedBody, r)
		fmt.Fprintf(w, `{
			"name":"test",
			"mongoDBMajorVersion":"3.4",
			"replicationFactor":3,
			"backupEnabled":false,
			"providerBackupEnabled":false,
			"diskSizeGB":5,
			"paused":false,
			"autoScaling":{
				"diskGBEnabled":false
			},
			"providerSettings":{
				"providerName":"AWS",
				"regionName":"US_EAST_1",
				"instanceSizeName":"M0"
			},
			"replicationSpec":{
				"US_EAST_1":{
					"priority":7,
					"electableNodes":2,
					"readOnlyNodes":1
				}
			}
		}`)
	})

	client := NewClient(httpClient)
	params := &Cluster{
		DiskSizeGB: float64(5),
	}
	cluster, _, err := client.Clusters.Update("123", "test", params)
	providerSettings := ProviderSettings{ProviderName: "AWS", RegionName: "US_EAST_1", InstanceSizeName: "M0"}
	replicationSpec := map[string]ReplicationSpec{
		"US_EAST_1": ReplicationSpec{Priority: 7, ElectableNodes: 2, ReadOnlyNodes: 1},
	}
	expected := &Cluster{
		Name:                "test",
		MongoDBMajorVersion: "3.4",
		ReplicationFactor:   3,
		BackupEnabled:       false,
		Paused:              false,
		DiskSizeGB:          float64(5),
		ProviderSettings:    providerSettings,
		ReplicationSpec:     replicationSpec,
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, cluster)
}

func TestClusterService_Delete(t *testing.T) {
	httpClient, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/api/atlas/v1.0/groups/123/clusters/test", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, "DELETE", r)
		fmt.Fprintf(w, `{}`)
	})

	client := NewClient(httpClient)
	resp, err := client.Clusters.Delete("123", "test")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
