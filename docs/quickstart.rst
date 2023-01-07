.. warning::

    The examples below are slightly outdated and will be revisited at some point.
    All commands should still work, but the output might be a little different now.
    Please refer to the :ref:`getting_started` guide for a more up-to-date version.

.. _quickstart:

Quickstart
==========

This does not really explain the philosophy behind ``floo``, but gives a good
idea what the tool is able to do and how it's supposed to be used. Users
familiar to ``git`` should be able to grok most of the commands intuitively.


1. Adding files
---------------

Before synchronizing them, you need to *stage* them. The files will be stored
encrypted (and possibly compressed) in blobs on your hard disks.

.. raw:: html

    <script id="asciicast-ZUnAzK1qdvvgjsmiL2VmmPEzX" src="https://asciinema.org/a/ZUnAzK1qdvvgjsmiL2VmmPEzX.js" async></script>


2. Coreutils
------------

``floo`` provides implementations of most file related core utils like ``mv``,
``cp``, ``rm``, ``mkdir`` or ``cat``. Handling of files should thus feel
familiar for users that know the command line.

.. raw:: html

    <script id="asciicast-pTdBvvgtM5p9b9MENRuNzMQga" src="https://asciinema.org/a/pTdBvvgtM5p9b9MENRuNzMQga.js" async></script>

3. Mounting
-----------

For daily use and for use with other tools you might prefer a folder that
contains the file you gave to ``floo``. This can be done via the built-in FUSE
layer.

.. raw:: html

    <script id="asciicast-D9YOEI77GqbTSgmv5UwGZlucX" src="https://asciinema.org/a/D9YOEI77GqbTSgmv5UwGZlucX.js" async></script>

.. note::

    Some built-in commands provided by floo are faster.
    ``floo cp`` for example only copies metadata, while the real ``cp`` will copy the whole file.

If you wish to always have the mount when ``floo`` is running, you should look
into :ref:`permanent-mounts`.

4. Commits
----------

In it's heart, ``floo`` is very similar to ``git`` and also supports versioning
via commits. In contrast to ``git`` however, there are no branches and you
can't go back in history -- you can only bring the history back up front.

.. raw:: html

   <script id="asciicast-4omRS5anJthenxWnpdLLxwNTv" src="https://asciinema.org/a/4omRS5anJthenxWnpdLLxwNTv.js" async></script>

