# SOPS: Secrets OPerationS

`sops` is a secrets management solution and editor of encrypted files that
supports a variety of file formats and encryption providers.

Here is a quick demo of `sops` in action:

![SOPS in action](https://i.imgur.com/X0TM5NI.gif)


A more in-depth overview is available as part of Julien Vehent's Securing
DevOps Show & Tell:

<iframe width="560" height="315" src="https://www.youtube-nocookie.com/embed/V2PRhxphH2w" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

If you want to use SOPS as a Go library, take a look at the [decrypt
package](https://godoc.org/go.mozilla.org/sops/decrypt).

**Questions?** Ping `ulfr` and `autrilla` in `#security` on
[`irc.mozilla.org`](https://wiki.mozilla.org/IRC) (you can use a web client
like [mibbit](https://chat.mibbit.com)).

**What happened to Python SOPS?** We rewrote `sops` in Go to solve a number of
deployment issues, but the Python branch still exists under `python-sops`. We
will keep maintaining it for a while, and you can still `pip install sops`, but
we strongly recommend you use the Go version instead.

Backward compatibility
----------------------

We strive to make as few backwards-incompatible changes as possible to the
`sops` command line tool. We follow [Semantic Versioning](https://semver.org/),
so in the rare occurence that we break compatibility on the CLI, you'll know.

The file format will always be backwards compatible: this means that newer
versions of SOPS will be able to load files created with older versions of
SOPS.

Security
--------

Please report security issues to jvehent at mozilla dot com, or by using one of
the contact method available on keybase: https://keybase.io/jvehent

License
-------

Mozilla Public License Version 2.0

Authors
-------

The core team is composed of:

* Adrian Utrilla [@autrilla](https://github.com/autrilla)
* Julien Vehent [@jvehent](https://github.com/jvehent)
* AJ Banhken [@ajvb](https://github.com/ajvb)

And a whole bunch of [contributors](https://github.com/mozilla/sops/graphs/contributors).

Credits
-------

SOPS was inspired by [hiera-eyaml](https://github.com/TomPoulton/hiera-eyaml),
[credstash](https://github.com/LuminalOSS/credstash),
[sneaker](https://github.com/codahale/sneaker), [password
store](http://www.passwordstore.org/), and too many years managing PGP
encrypted files by hand.
