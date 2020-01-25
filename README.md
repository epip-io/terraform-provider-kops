# Terraform Provider for Kops

## Manual Install

To manually install, first run `go generate`, this will create a vendor directory for the dependences as well as generating a lot of the provider code. Further, it will copy `./hack/bindata.go` (see [Notes](#Notes)) into `./vendor/k8s.io/kops/upup/models`. This is required to be able to build the provider.

Terraform providers are installed into the users [Terraform plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) directory. So after generating code, so to build and install run:

```bash
mkdir -p ~/terraform.d/plugins/
go build -mod vendor -o ~/.terraform.d/plugins/$(basename $(pwd))_$(git describe --exact-match --match 'v*')
```

## Notes

- `go.mod` uses a lot of `replace` lines that are taken from the github.com/kubernetes/kops commit that aligns with the version used to generate the code. This assures that the provider uses the same underlying Kubernetes module versions.
- `bindata.go` is generated from github.com/kubernetes/kops/upup/models, using the following command, and then copied into `./hack`:
  
  ```go-bindata -pkg models -prefix $(pwd) -prefix upup/models cloudup/... nodeup/...```

  (`go-bindata` is gotten using `go get -u github.com/jteeuwen/go-bindata/go-bindata`)

## Links

- [Kops API Spec](https://kops.sigs.k8s.io/cluster_spec/)
- [Terraform providers](https://www.terraform.io/docs/configuration/providers.html)
- [Terraform custom providers](https://www.terraform.io/docs/extend/writing-custom-providers.html)

## To-Do

- [x] From Kops Go Cluster and InstanceGroup `struct`, generate code for
  - [x] the Terraform provider schema
  - [x] Marshalling between Terraform and Kops schemas
- [x] Generate Terraform provider, and initial datasources and resources
- [ ] Generate testing
  - [ ] Terraform provider schema
  - [ ] Marshalling being schema
  - [ ] Provider
- [ ] Testing
  - [ ] `./util/templater`
- [ ] Implement CI/CD pipeline (Travis? CircleCI?)
- [ ] Determine method for automatically updating `./hack/bindata.go`
- [ ] Refactor `./hack/gen-go.go` to make it reusable
- [ ] Support for multiple Kops API versions
- [ ] Release binaries
