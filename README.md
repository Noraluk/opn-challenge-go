# GO-TAMBOON ไปทำบุญ

This is a small challenge project to see how good you are with Go. Included in this
repository is a CSV list of Song-pah-pa (ซองผ้าป่า). It is a white envelope with money
inside and the donor name printed on the front. They are usually collected en bulk from
multiple people in order to round up money to repair or construct new temple buildings.

The idea is that your donation amount should be kept secret les the activity becomes an
act of flaunting your wealth.

But we're a payment gateway, we can do better than that. The envelope will contain,
instead, a valid CC number (fake ones, not a real working card) and the desired donation
amount. The entire list is also encrypted using NSA-proof variant of the
[Caesar Cipher][1] :troll:

### CONTENTS

- `data/fng.csv.rot128` - A ROT-128 encrypted CSV file.
- `cipher/rot128.go` - Sample ROT-128 encrypt-/decrypter.

### EXERCISE

Write a GO command-line module that, when given the CSV list, calls the [Charge API][0] to
make donations by creating a charge for each row in the file and produce a summary at the
end.

Example:

```sh
$ cd $GOPATH/omise/go-tamboon
$ go install -v .

$ $GOPATH/bin/go-tamboon test.csv

performing donations...
done.

        total received: THB  210,000.00
  successfully donated: THB  200,000.00
       faulty donation: THB   10,000.00

    average per person: THB      534.23
            top donors: Obi-wan Kenobi
                        Luke Skywalker
                        Kylo Ren
```

**Requirements:**

- Decrypt the file using a simple [ROT-128][2] algorithm.
- Make donations by creating a Charge via the [Charge API][0] for each row in the
  decrypted CSV.
- Produce a brief summary at the end.
- Handle errors gracefully, without stopping the entire process.
- Writes readable and maintainable code.

**Bonus:**

- Have a good Go package structure.
- Be a good internet citizen and throttles the API call if we hit rate limit.
- Run as fast as possible on a multi-core CPU.
- Allocate as little memory as possible.
- Complete the entire process without leaving large trace of Credit Card numbers
  in memory, or on disk.
- Ensure reproducible builds on your workspace.

[0]: https://www.omise.co/charges-api
[1]: https://en.wikipedia.org/wiki/Caesar_cipher
[2]: https://play.golang.org/p/dCWYyWPHwj4

### How to run it

1. Please fill the public key and secret key of OPN into config/config.yaml
2. You can run by:

```
./go-tamboon ./data/fng.1000.csv.rot128
```

or

```
go run main.go ./data/fng.1000.csv.rot128
```

### FYI

1. Because I use credential unpaid user when I run the code I found the 403 response after I run customers more than 60-70 ( I predict ) from the omise which means it is the limitation of the unpaid user.
2. I set the default log level to info level, so If you want to see log you can change the log level in config.yaml to debug.
