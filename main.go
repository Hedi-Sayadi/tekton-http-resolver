package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/resolution/common"
	"github.com/tektoncd/pipeline/pkg/resolution/resolver/framework"
	"knative.dev/pkg/injection/sharedmain"
)

type HTTPResolver struct{}

func main() {
	log.Println("Starting Tekton Resolver...")
	sharedmain.Main("controller",
		framework.NewController(context.Background(), &HTTPResolver{}),
	)
	log.Println("Tekton Resolver started.")

}

// Initialize initializes the resolver.
func (r *HTTPResolver) Initialize(ctx context.Context) error {
	// Add initialization logic here if needed
	return nil
}

// GetName returns the name of the resolver.
func (r *HTTPResolver) GetName(context.Context) string {
	return "HTTPResolver"
}

// GetSelector returns a label selector for the resolver.
func (r *HTTPResolver) GetSelector(context.Context) map[string]string {
	return map[string]string{
		common.LabelKeyResolverType: "http",
	}
}

// ValidateParams validates the parameters passed to the resolver.
func (r *HTTPResolver) ValidateParams(ctx context.Context, params []pipelinev1.Param) error {
	// Add validation logic here if needed
	if len(params) > 0 {
		return errors.New("no params allowed")
	}
	return nil
}

// Resolve fetches CRDs from HTTP URLs.
func (r *HTTPResolver) Resolve(ctx context.Context, params []pipelinev1.Param) (framework.ResolvedResource, error) {
	sourceURL := os.Getenv("RESOLVER_SOURCE_URL")
	if sourceURL == "" {
		return nil, errors.New("RESOLVER_SOURCE_URL must be configured")
	}

	crd, err := fetchCRD(ctx, sourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resource from %s: %w", sourceURL, err)
	}

	return &CRDResource{data: crd}, nil
}

// CRDResource represents a resolved CRD fetched from an HTTP URL.
type CRDResource struct {
	data []byte
}

// Data returns the CRD data.
func (c *CRDResource) Data() []byte {
	return c.data
}

// Annotations returns any metadata needed alongside the CRD data.
func (c *CRDResource) Annotations() map[string]string {
	return nil
}

// RefSource returns the source reference of the remote data.
func (c *CRDResource) RefSource() *pipelinev1.RefSource {
	return nil
}

// // fetchCRD fetches the CRD from the given URL.
// func fetchCRD(url string) ([]byte, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

//		return body, nil
//	}
func fetchCRD(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error fetching CRD from %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from %s: %v", url, err)
		return nil, err
	}

	return body, nil
}
