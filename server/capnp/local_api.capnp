using Go = import "/go.capnp";
using User = import "../../gateway/db/capnp/user.capnp";

$Go.Package("capnp")
$Go.import("/server/capnp")

struct Hint {
    path            @0 :Text;
    encryptionAlgo  @1 :Text;
    compressionAlgo @2 :Text;
}