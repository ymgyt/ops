===
ops
===

ops is utility tool kit for daily operations.


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