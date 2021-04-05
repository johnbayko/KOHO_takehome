I did this assignment for a position at KOHO. I didn't get it, but decided to
put it up as a programming example.

The review said the main problem was a lack of documentation. I wasn't sure
what do document, since the design seemed "obvious" to me. I knew the next step
would be a technical interview where I could answer questions and get a better
idea of what they wanted in that area, but it seems understanding that was part
of the exercise, so I didn't get to the interview part.

Since then I looked at some versions done by others on github, and have a
better idea of other approaches, and how to explain my design.

The review also said it could be more concise. I generally emphasize some
things over conciseness, in this case:

- Readability.

- Modifiability.

- Bug resistence. I usually rank this higher for languagers like C and C++, or
  Java, but golang is a bit limited in language features like local constants,
  smart enumerations, and the like while also being less prone to bugs to begin
  with, so clarity becomes more important.

Requirements
============

The requirements were laid out in the README.md file.

Environment
===========

The assignmet had no details of the environment, other than it was file based,
with a chronological input file and an output file.

In a real application, even if input is in batches like the example, the
transactions and accounts would need to persist, so a database would be needed.
I was surprised that only two other solutions mentioned or implemented
persistent storage.

Peristence can be accomplished many ways, so the interface to the program
should be independent of the implementation.

It's also likely that the loading would be made available as a web service or
other single transaction interface, so the batch file handling should be
separate from handling individual transactions.

Strategy
========

The system is broken up into separate packages:

- The main interface for the current assignment, which reads the input, writes
  the output, and reports errors.

- The main handling package which applies the rules for the transaction based
  on information from storage, indicates to the calling package whether it's
  accepted or rejected, and passes it on to the storage for updating.

- The storage interface, which defines the methods to query and update the
  transaction history, and current accounts.

- The storage implementation.

Implementation
==============

Main updater
------------

The ``main`` routine is separate from the ``update`` function for two reasons.

The main routine doesn't return, it calls ``os.exit()`` to allow an error code
to be returned as an exit status. That means ``defer`` statements don't get
run, so files won't be closed. It might not be necessary to close a file on
exit, but better to be sure (buffered output can be a problem in some systems).

The main routine constructs the handler and storage objects and pass them to
the update function, allowing them to be customised without any changes to the
update function. The basic golang unit testing doesn't support mocking
functions and objects, so in general it's better to explicitly pass
dependencies as parameters.

The update routine then opens the input and output files (with default names if
not specified). The assignment said I can assume that input is one transaction
per line, but the input stream can be passed to create a json ``Decoder``,
which will decode records regardless of how many lines they're on, and the
code is simpler. Similarly an ``Encoder`` object is created from the output
file.

Some of the input fields need parsing, so ``getTransaction()`` encapsulates all
code to read a transaction and parse any fields to produce a transaction record
the rest of the program uses.

Similarly, ``putAcceptance()`` constructs and writes the output record from the
original transaction record.

Between the two, the transaction is passed to the update handler ``Load()``
method, which returns a flag indicating acceptance or rejection, and an error.

The assignment said to ignore duplicate transactions. A duplicate is indicated
by a specific error return, which indcates not to output an error message. It
could have been done with an additional flag, but since an error object is
returned anyway, this simplifies the return value list. A duplicate will never
be detected if there's another error, so there is no possible conflict.

The requirement to ignore duplicates could change, in which case the code
change is small.

Testing
~~~~~~~

There wasn't much to test at this point, but the parsing of the load amounts
could be tricky, so I added unit tests for that.

Transaction handling
--------------------

The transaction loading has three steps:

- Check for validation.

- Save the transaction and validation.

- Update the balance.

The assignment didn't specify that the customer balance should be saved, that
was just an assumption, and could be removed if the purpose is actually just
validation. In addition, I created a customer table to tie the transactions and
customer balance together, again as a placeholder that can be removed.

Validation is a separate method, which calls each validatiaon function in turn,
and exits if any fail. Individual validations are in separate functions because
there is some date manipulation needed that can be encapsulated (to avoid
duplication and limit the scope of any changes).

A side note, I noticed most examples use a loop to find the previous Monday. I
used a switch because that usually compiles down to jump table and single
executed instruction. Also it's very clear what it's doing. Even though there
is more source code, the compiled code is faster and still small.

The alternative would have been to cast ``time.Duration`` to ``int`` to do math
on them. One example did this. It would be efficient with less code, but I
don't like the idea of casts to subvert type protections. In this case it would
be justified because it's a problem with the golang ``time`` library design,
but I decided to work within golang as implemented.

Some rules were not specified in the requirements, which needed some trial and
error to figure out:

- Transaction ids are not universal, they are tied to customers, so several
  customers can have transactions with the same id.

- A rejected transaction is logged anyway. A later correct transaction with the
  same id must be skipped as a duplicate.

Storage
-------

API
~~~

The storage API was defined by the transaction handling implementation.
Basically the needs were:

- Get amount loaded for a period.

- Get transaction count for a period.

- Log the transaction.

- Update the customer balance.

Sample implementation
~~~~~~~~~~~~~~~~~~~~~

I decided ``sqlite3`` would be good for a sample implementation, since it
doesn't need additional infrastructure, is widely available and easy to install
(often by default), and still allows a proper database model.

This meant adding a dependency on a golang sqlite driver. Apart from that, the
implementation used the standard golang sql library.

The other implementation that used a database used mysql.

Testing
-------

For testing, a mock customer storage object was defined in the test file which
would check the method parameters and return test results. The test function
allocated a handler object with the mock storage object as a parameter to the
construction function.

A list of handler transactions and responses, and storage mock parameters and
reponses allow testing for all possible input parameters and storage
configurations. It was a little ambitious and I was only able to get one test
scenario coded for the submission deadline, with a bug in the mock
implementation. I fixed that in the current version.

Setup and execution
===================

This part of the documentation was included in the email.

``sqlite3`` is required to be installed. Apart from that, there are two scripts
to handle everything else.

``setup`` gets the golang sqlite driver, and creates the initial sqlite
database tables.

``testrun`` runs the unit tests, clears the databse, runs the actual update
program, and then compares the output results to the sample output.

