### mitum-nft-v2

*mitum-nft-v2* is a nft contract model based on the second version of mitum(aka mitum2).

#### Features,

* user defined contract account state policy: collection.
* collection: user defined nft collection.
* *LevelDB*, *Redis*: as mitum2 does, *LevelDB* and *Redis* can be primary storage.
* reference nft standard: ERC-721
* multiple collection policy for one contract account.

#### Installation

Before you build `mitum-nft`, make sure to run `docker run`.

```sh
$ git clone https://github.com/ProtoconNet/mitum-nft-v2

$ cd mitum-nft-v2

$ go build -o ./mitum-nft
```

#### Run

```sh
$ ./mitum-nft init <config file>

$ ./mitum-nft run <config file>
```

[standalong.yml](standalone.yml) is a sample of `config file`.
