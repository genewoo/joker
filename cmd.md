# Joker Tool Requirement

## General

Write a command line tool for poker player to calculate different types of calculation, it's git-like command line.
The command line tool will follow a git-like structure, where the main command is followed by a <gametype> and a <subcommand>. The tool will also support various options that can be passed to the subcommands.



## commandline formation

```bash
joker <gametype> <subcommand> <options>
```

## Gametypes
standard: A standard card game.
holdem: A Texas Hold'em style game.

## gametype : standard

It's a stnadard card games. After running the command, it will prompt for task for the next step.

### subcommand

```csv
deal
```

### options

```csv
-p, --players: Number of players in the game (default: 2).
-d, --decks: Number of decks to use (default: 1).
-j, --joker: Whether the deck includes jokers (default: false).
-k, --keep: Number of cards to keep (default: 0)
-n, --numberofcards number: Number of cards distributed to each player (default: calculated based on the number of cards in the deck and players).
```

<!-- ### tasks

- deal: Shuffle and deal cards to players.
- show N: Show cards of the Nth player (for standard gametype).
- eq: Calculate equity for players (for holdem gametype).
 -->

## gametype : holdem

It's a holdem game

### subcommand

```csv
deal
eq
```

### options

```csv
-p, --players number, how many players in the game (default: 2), p
-n, --numberofcards number, how many cards are distributed to each player (default: 2), n
-t, --type, enum, from values from texas, omaha, short, (default: texas)

```