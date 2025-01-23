# Rules

- 4 players, team up with 2, sitting opposite each other. Call Team A and Team B, A(1, 3) and B(2, 4)
- 2 decks of playing cards, jokers included (108 cards total)
- The initial status will be:
  -  Team A and B's current level, from 2 to A
  -  Indicate the last game's player rank, for example, 1,3,2,4

- From the initial status, the program will indicate following status
  - the dealer will be the last game's first player, take the previous example, 1
  - dealing cards from the last game's last player to the rest of players, take the previous example, 4,1,2,3
  - The game current level is the level from the winner team.
- Card ranking are noraml as Default, from 2 to A, and jokers are the biggest. If the game level is 5, then 5 is the bigger than A, and 2 is the smallest.
- After dealing there are swap card rules:
  - If the last game, the last two players are both from the same team, the last game's last two player should give the biggest card to the last game's first two player, except the loser team's players got two red jokers. If the last game's last two players are from different team, the last game's last player should give the biggest card to the last game's first player, except the player got two red jokers.
  - The card to be given should be the biggest card in the hand, and it should not be the level card and suit in Hearts.
  - As a return, any player recieved the card should give back any card to the giver, it should not return card bigger than 10.
- If the swapping happens, the dealer will be the player gave the biggest card to the first player.