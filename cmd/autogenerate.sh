# Copyright 2017 caicloud authors. All rights reserved.

# The script auto-generate kubernetes clients, listers, informers

set -e

# Add your packages here
PKGS=(
  release/v1alpha1
  config/v1alpha1
  apiextensions/v1beta1
  loadbalance/v1alpha2
  resource/v1alpha1
)

CLIENT_PATH=github.com/caicloud/clientset
CLIENT_APIS=$CLIENT_PATH/pkg/apis
GO_HEADER_FILE="../hack/boilerplate/boilerplate.go.txt"

for path in "${PKGS[@]}"
do
	ALL_PKGS="$CLIENT_APIS/$path "$ALL_PKGS
done

function codegen::join() { local IFS="$1"; shift; echo "$*"; }

function join {
	local IFS="$1"
   	shift
   	result="$*"
}

join "," ${PKGS[@]}
PKGS=$result

join "," $ALL_PKGS
FULL_PKGS=$result

echo "PKGS: $PKGS"
echo "FULL PKGS: $FULL_PKGS"

cd $(dirname ${BASH_SOURCE[0]})

${GOPATH}/bin/deepcopy-gen --input-dirs $(codegen::join , "$FULL_PKGS") \
-O zz_generated.deepcopy --bounding-dirs "$CLIENT_APIS" --go-header-file "$GO_HEADER_FILE"

echo "Generated deepcopy funcs"

go run ./client-gen/main.go \
  -n kubernetes \
  --clientset-path $CLIENT_PATH \
  --input-base $CLIENT_APIS \
  --input $PKGS

echo "Generated clients"

go run ./lister-gen/main.go \
  -p $CLIENT_PATH/listers \
  -i $FULL_PKGS

echo "Generated listers"

go run ./informer-gen/main.go \
  -p $CLIENT_PATH/informers \
  --versioned-clientset-package $CLIENT_PATH/kubernetes \
  --listers-package $CLIENT_PATH/listers \
  -i $FULL_PKGS

echo "Generated informers"

cd - >/dev/null
