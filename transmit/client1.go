package transmit

import (
	"fmt"

	"github.com/TremblingV5/TinyCache/pb"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpClient struct {
	baseURL string
}

func (h *httpClient) Get(in *pb.GetKeyRequest, out *pb.GetKeyResponse) error {
	u := fmt.Sprintf(
		"%v/%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetBucket()),
		url.QueryEscape(in.GetKey()),
	)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

func (h *httpClient) Set(in *pb.SetKeyRequest, out *pb.SetKeyResponse) error {

}

func (h *httpClient) Delete(in *pb.DeleteKeyRequest, out *pb.DeleteKeyRequest) error {
	//TODO implement me
	panic("implement me")
}

func (h *httpClient) DeleteBucket(in *pb.DeleteBucketRequest, out *pb.DeleteBucketResponse) error {
	//TODO implement me
	panic("implement me")
}

var _ tinycache.PeerHandle = (*httpClient)(nil)
