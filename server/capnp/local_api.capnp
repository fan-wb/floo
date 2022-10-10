using Go = import "/go.capnp";
using User = import "../../gateway/db/capnp/user.capnp";

@0xea883e7d5248d81b;
$Go.package("capnp");
$Go.import("/server/capnp/local_api");

# nnn@hal9000:~/code/floo$ capnp compile -I$GOPATH/src/capnproto.org/go/capnp/std -ogo server/capnp/local_api.capnp


struct Hint {
    path            @0 :Text;
    encryptionAlgo  @1 :Text;
    compressionAlgo @2 :Text;
}