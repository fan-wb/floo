using Go = import "/go.capnp";

@0xa0b1c18bd0f965c4;

$Go.package("capnp");
$Go.import("/gateway/db/capnp/user");

struct User {
	name         @0 :Text;
	passwordHash @1 :Text;
	salt         @2 :Text;
	folders      @3 :List(Text);
	rights       @4 :List(Text);
}