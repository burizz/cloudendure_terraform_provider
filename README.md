# Terraform Provider cloudendure

Run the following command to build the provider

```shell
export OS_ARCH="$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
#go build -o terraform-provider-cloudendure

go build -gcflags="all=-N -l" -o terraform-provider-cloudendure

mkdir -p ~/.terraform.d/plugins/hashicorp.com/edu/cloudendure/0.2/$OS_ARCH

mv terraform-provider-cloudendure ~/.terraform.d/plugins/hashicorp.com/edu/cloudendure/0.2/$OS_ARCH

rm ./examples/.terraform.lock.hcl

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```


TF_REATTACH_PROVIDERS='{"registry.terraform.io/hashicorp.com/edu/cloudendure":{"Protocol":"grpc","Pid":3382870,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin713096927"}}}'

dlv exec --headless ./terraform-provider-cloudendure -- --debug


ps | grep cloudendure