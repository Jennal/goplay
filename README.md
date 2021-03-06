# GoPlay

[![GoDoc](https://godoc.org/github.com/Jennal/goplay?status.svg)](https://godoc.org/github.com/Jennal/goplay) [![Go Report Card](https://goreportcard.com/badge/github.com/jennal/goplay)](https://goreportcard.com/report/github.com/jennal/goplay)

GoPlay is a framework for game service, written by pure golang.

## Clients

- Unity3d/C#: [GoPlay-Client-Unity3d](https://github.com/Jennal/goplay-client-unity3d)
- C++: [GoPlay-Client-Cpp](https://github.com/Jennal/goplay-client-cpp)
- Javascript: [GoPlay-Client-Javascript](https://github.com/Jennal/goplay-client-javascript)

## Cluster Services

- [Master](https://github.com/Jennal/goplay-master): master manages all servers, and make service cluster
- [Connector](https://github.com/Jennal/goplay-connector): connector is gate way of client

### Run Cluster Services By Docker

```docker
docker run --name goplay-master --rm -i -t -p 6812:6812 jennal/goplay-master
docker run --name goplay-connector --rm -i -t -p 9934:9934 --link goplay-master jennal/goplay-connector --master-host goplay-master
```

## Demos

- [GoPlay-Demos](https://github.com/Jennal/goplay-demos)

## The MIT License

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.