===
ops
===

ops is utility tool kit for daily operations.

Install
=======

.. code-block:: bash

   $ curl -fsSL -o ops https://github.com/ymgyt/ops/releases/download/v0.0.2/ops-linux-amd64-v0.0.2
   $ chmod +x ops
   $ ./ops version
   v0.0.2


Usage
=====

du
---

.. code-block:: bash

   du print provided directory disk usage recursively.

   USAGE
     ops du [OPTIONS] <root_directory>

   Options
     -h, --help         : print this.
         --humanize     : humanize bytes. ex MB
     -i, --ibytes       : humanize ibytes. ex MiB
     -l, --level        : print directory level.
         --max-recursion: max recursion.
         --order        : output disk usaage order.
     -r, --relative     : print relative file path from root.
     -v, --verbose      : verbose


rand
----

.. code-block:: bash

   rand generate random characters.

   USAGE
    ops rand [OPTIONS]

   default character pool: abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789


   Options
     -h, --help         : print this.
         --len          : random strings length.
         --special-chars: characters used in addition to the standard string pool for random string generation



.. code-block:: bash

   ops rand --len=50 --special-chars='?!_-*.,&%'
   XyACZMsL?tPJ_Vd8Yg93!EGoFpqOuRDHWV*jFP*r7u2OqIP4C_