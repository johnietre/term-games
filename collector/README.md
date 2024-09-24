# collector
A game to see who can collect the most balls/dots. The balls are randomly spawned in the play area and players compete to get to them. Multiple players can occupy a space at the same time, but should a ball spawn on that space, the first player there gets the ball (or maybe nothing can spawn there?).
The creator of a game specifies the dimensions of the play area and can specify whether each player must conform to said dimensions. If conformity is required, players that do not conform are not allowed in a newly created game. Players the lose conformity are temporarily blocked from moving until returning to appropriate dimensions.

## Implementation Ideas
- The server holds as little game state as possible, and manages state by synchronizing using the players (majority is truth)
- Server holds all game state
