# Blackjack Simulator

I've always been fascinated by Blackjack. Some pros say that, if you follow a
basic strategy, your odds of winning go up significantly. So, that got me
wondering:

> If you follow a Blackjack strategy algorithmically, how will you do?

This tiny app is meant to address that.

## How does it work?

You author a strategy with a pretty straight forward DSL. The app will run this
strategy against a given number of games (default 100) and output how it does.
Here's an example strategy.

```
[soft]
   2 3 4 5 6 7 8 9 10 A
13 H H H H H H H H  H H
14 H H H H H H H H  H H
15 H H H H H H H H  H H
16 H H H H H H H H  H H
17 S S S S S S S S  S S
18 S S S S S S S S  S S
19 S S S S S S S S  S S
20 S S S S S S S S  S S
21 S S S S S S S S  S S

[hard]
   2 3 4 5 6 7 8 9 10 A
 4 H H H H H H H H  H H
 5 H H H H H H H H  H H
 6 H H H H H H H H  H H
 7 H H H H H H H H  H H
 8 H H H H H H H H  H H
 9 H H H H H H H H  H H
10 H H H H H H H H  H H
11 H H H H H H H H  H H
12 H H H H H H H H  H H
13 H H H H H H H H  H H
14 H H H H H H H H  H H
15 H H H H H H H H  H H
16 H H H H H H H H  H H
17 S S S S S S S S  S S
18 S S S S S S S S  S S
19 S S S S S S S S  S S
20 S S S S S S S S  S S
21 S S S S S S S S  S S
```

The `[soft]` section describes soft-hand strategy. The `[hard]` section
describes hard-hand strategy.

You can run that strategy through the simulator like this.

```
$ ./go-blackjack --strategy=strategies/passive --games=10000
2013/11/09 22:31:07 Loading strategy strategies/passive
2013/11/09 22:31:09 Total Hands         551588
2013/11/09 22:31:09 Total Wins          213924  (38.783%)
2013/11/09 22:31:09 Total Losses        277828  (50.369%)
2013/11/09 22:31:09 Total Pushes        49836   (9.035%)
```

## Does it actually work??

I dunno. So far, I've tried two different strategies and here are my results for each.

### Passive Strategy

This strategy is checked in to the repo.

```
$ ./go-blackjack --strategy=strategies/passive --games=100000
2013/11/09 22:32:12 Loading strategy strategies/passive
2013/11/09 22:32:33 Total Hands         5515165
2013/11/09 22:32:33 Total Wins          2141896 (38.836%)
2013/11/09 22:32:33 Total Losses        2780783 (50.421%)
2013/11/09 22:32:33 Total Pushes        492486  (8.930%)
```

### Wizard of Odds Strategy

This strategy is also checked in to the repo and described on the [Wizard of
Odds](http://wizardofodds.com/games/blackjack/) website.

**NOTE:** One big missing piece that is described in the Wizard of Odds
strategy that is missing here is splitting. This simulator does not support it!

```
$ ./go-blackjack -strategy strategies/wizard_simple -games 100000
2013/11/09 22:25:13 Loading strategy strategies/wizard_simple
2013/11/09 22:25:33 Total Hands       5562401
2013/11/09 22:25:33 Total Wins        2275528 (40.909%)
2013/11/09 22:25:33 Total Losses      2761467 (49.645%)
2013/11/09 22:25:33 Total Pushes      425406  (7.648%)
```

## Assumptions

1. `/dev/urandom` is sufficiently random for our purposes.
1. Not shuffling between hands is OK.
1. Simulates a six-deck shoe by default.

## Contributing

You know what do! Fork and submit a pull request. Strategies are, of course,
welcome as well.
