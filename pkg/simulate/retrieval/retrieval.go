package retrieval

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ethersphere/beekeeper/pkg/bee"
	"github.com/ethersphere/beekeeper/pkg/beeclient/api"
	"github.com/ethersphere/beekeeper/pkg/beekeeper"
	"github.com/ethersphere/beekeeper/pkg/random"
	"github.com/prometheus/client_golang/prometheus/push"
)

// Options represents simulation options
type Options struct {
	ChunksPerNode   int // number of chunks to upload per node
	GasPrice        string
	MetricsPusher   *push.Pusher
	PostageAmount   int64
	PostageDepth    uint64
	PostageLabel    string
	PostageWait     time.Duration
	Seed            int64
	UploadNodeCount int
	UploadDelay     time.Duration
}

// NewDefaultOptions returns new default options
func NewDefaultOptions() Options {
	return Options{
		ChunksPerNode:   1,
		GasPrice:        "",
		MetricsPusher:   nil,
		PostageAmount:   1000,
		PostageDepth:    16,
		PostageLabel:    "test-label",
		PostageWait:     5 * time.Second,
		Seed:            random.Int64(),
		UploadNodeCount: 1,
		UploadDelay:     5 * time.Second,
	}
}

// compile simulation whether Upload implements interface
var _ beekeeper.Action = (*Simulation)(nil)

// Simulation instance
type Simulation struct{}

// NewSimulation returns new upload simulation
func NewSimulation() beekeeper.Action {
	return &Simulation{}
}

// Run executes retrieval simulation
func (s *Simulation) Run(ctx context.Context, cluster *bee.Cluster, opts interface{}) (err error) {
	o, ok := opts.(Options)
	if !ok {
		return fmt.Errorf("invalid options type")
	}

	rnds := random.PseudoGenerators(o.Seed, o.UploadNodeCount)
	fmt.Printf("Seed: %d\n", o.Seed)

	metrics := newMetrics(cluster.Name()+"-"+time.Now().UTC().Format("2006-01-02-15-04-05-000000000"), o.MetricsPusher)

	overlays, err := cluster.FlattenOverlays(ctx)
	if err != nil {
		return err
	}
	clients, err := cluster.NodesClients(ctx)
	if err != nil {
		return err
	}

	// continually upload chunk and download
	for {
		sortedNodes := cluster.NodeNames()
		for i := 0; i < o.UploadNodeCount; i++ {

			nodeName := sortedNodes[i]
			client := clients[nodeName]

			batchID, err := client.GetOrCreateBatch(ctx, o.PostageAmount, o.PostageDepth, o.GasPrice, o.PostageLabel)
			if err != nil {
				fmt.Printf("error: node %s: batch id %v\n", nodeName, err)
				continue
			}
			fmt.Printf("node %s: batch id %s\n", nodeName, batchID)
			time.Sleep(o.PostageWait)

			for j := 0; j < o.ChunksPerNode; j++ {
				chunk, err := bee.NewRandomChunk(rnds[i])
				if err != nil {
					fmt.Printf("error: node %s: %v\n", nodeName, err)
					continue
				}

				tag, err := client.CreateTag(ctx)
				if err != nil {
					return fmt.Errorf("create tag on node %s: %w", nodeName, err)
				}

				// upload chunk
				t0 := time.Now()
				ref, err := client.UploadChunk(ctx, chunk.Data(), api.UploadOptions{
					BatchID: batchID,
					Tag:     tag.Uid,
				})
				d0 := time.Since(t0)
				if err != nil {
					metrics.notUploadedCounter.WithLabelValues(overlays[nodeName].String()).Inc()
					if o.MetricsPusher != nil {
						if err := o.MetricsPusher.Push(); err != nil {
							fmt.Printf("upload node %s: %v\n", nodeName, err)
						}
					}
					fmt.Printf("error: node %s: %v\n", nodeName, err)
					continue
				}
				metrics.uploadedCounter.WithLabelValues(overlays[nodeName].String()).Inc()
				metrics.uploadTimeGauge.WithLabelValues(overlays[nodeName].String(), ref.String()).Set(d0.Seconds())
				metrics.uploadTimeHistogram.Observe(d0.Seconds())
				fmt.Printf("Chunk %s uploaded successfully to node %s\n", chunk.Address().String(), overlays[nodeName].String())

				// check if chunk is synced
				t1 := time.Now()
				err = client.WaitSync(ctx, tag.Uid)
				d1 := time.Since(t1)
				if err != nil {
					metrics.notSyncedCounter.WithLabelValues(overlays[nodeName].String()).Inc()
					if o.MetricsPusher != nil {
						if err := o.MetricsPusher.Push(); err != nil {
							fmt.Printf("upload node %s: %v\n", nodeName, err)
						}
					}
					fmt.Printf("sync with node %s: %v\n", nodeName, err)
					continue
				}
				metrics.syncedCounter.WithLabelValues(overlays[nodeName].String()).Inc()
				metrics.syncTagsTimeGauge.WithLabelValues(overlays[nodeName].String(), ref.String()).Set(d1.Seconds())
				metrics.syncTagsTimeHistogram.Observe(d1.Seconds())
				fmt.Printf("Chunk %s synced successfully with node %s\n", chunk.Address().String(), nodeName)

				// pick a random node to validate that the chunk is retrievable
				downloadNode := sortedNodes[rnds[i].Intn(len(sortedNodes))]

				// download chunk
				t2 := time.Now()
				data, err := clients[downloadNode].DownloadChunk(ctx, ref, "")
				d2 := time.Since(t2)
				if err != nil {
					metrics.notDownloadedCounter.WithLabelValues(overlays[downloadNode].String()).Inc()
					if o.MetricsPusher != nil {
						if err := o.MetricsPusher.Push(); err != nil {
							fmt.Printf("upload node %s, download node %s: %v\n", nodeName, downloadNode, err)
						}
					}
					fmt.Printf("error: node %s: %v\n", downloadNode, err)
					continue
				}
				metrics.downloadedCounter.WithLabelValues(overlays[downloadNode].String()).Inc()
				metrics.downloadTimeGauge.WithLabelValues(overlays[downloadNode].String(), ref.String()).Set(d2.Seconds())
				metrics.downloadTimeHistogram.Observe(d2.Seconds())

				// validate that chunk is retrieved correctly
				if !bytes.Equal(chunk.Data(), data) {
					metrics.notRetrievedCounter.WithLabelValues(overlays[downloadNode].String()).Inc()
					if o.MetricsPusher != nil {
						if err := o.MetricsPusher.Push(); err != nil {
							fmt.Printf("upload node %s, download node %s: %v\n", nodeName, downloadNode, err)
						}
					}
					fmt.Printf("Node %s. Chunk %d not retrieved successfully. Uploaded size: %d Downloaded size: %d Node: %s Chunk: %s\n", downloadNode, j, chunk.Size(), len(data), overlays[downloadNode].String(), ref.String())
					if bytes.Contains(chunk.Data(), data) {
						fmt.Printf("Downloaded data is subset of the uploaded data\n")
					}
					continue
				}
				metrics.retrievedCounter.WithLabelValues(overlays[downloadNode].String()).Inc()
				fmt.Printf("Chunk %s retrieved successfully from node %s\n", chunk.Address().String(), overlays[downloadNode].String())

				if o.MetricsPusher != nil {
					if err := o.MetricsPusher.Push(); err != nil {
						fmt.Printf("upload node %s, download node %s: %v\n", nodeName, downloadNode, err)
					}
				}
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(o.UploadDelay):
		}
	}
}
