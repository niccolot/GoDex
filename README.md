# GoDex

Pokemon-like CLI game in which you can explore the pokemon world, battle pokemons and collect them in your personal pokedex.

All the data are fetched from the [PokeAPI](https://pokeapi.co/) and subsequently cached for better performance in case of bad interent connection.

## How to use 

```
git clone https://github.com/niccolot/GoDex
cd GoDex
go build && ./GoDex 
```

### Requirements

* Go 1.22.5
* Internet connection

## How to play

Once started the program, a REPL prompt is initiated

```
Pokedex >
```

Enter the command `help` to see the possible actions you can take. 

Enter the command `exit` to quit the game.

Use up and down arrows to navigate through previosly used commands.

In order to go around the world enter the command `map`, which fetches the nearest 20 locations from the [PokeAPI](https://pokeapi.co/). 

```
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
sunyshore-city-area
sinnoh-pokemon-league-area
oreburgh-mine-1f
oreburgh-mine-b1f
valley-windworks-area
eterna-forest-area
fuego-ironworks-area
mt-coronet-1f-route-207
mt-coronet-2f
mt-coronet-3f
mt-coronet-exterior-snowfall
mt-coronet-exterior-blizzard
mt-coronet-4f
mt-coronet-4f-small-room
mt-coronet-5f
mt-coronet-6f
mt-coronet-1f-from-exterior
```

Enter the command multiple times to explore farther locations. If you want to go back use the command `mapb`.

In order to explore some areas use the command `explore` to list the local pokemons

```
Pokedex > explore canalave-city-area
-  tentacool
-  tentacruel
-  staryu
-  magikarp
-  gyarados
-  wingull
-  pelipper
-  shellos
-  gastrodon
-  finneon
-  lumineon
```

When you are in the vicinity of a pokemon you can try to cath it with the `catch` command

```
Pokedex > catch tentacool
Throwing a pokeball at tentacool...
tentacool was caugth and added to the pokedex!
```

If you fail capturing a pokemon it will escape and hide for a while.

When you have succesfully catch a pokemon you can find it listed in your pokedex, from which you can list it´s stats

```
Pokedex > pokedex
- tentacool


Pokedex > inspect tentacool
Name: tentacool
Height: 9
Weight: 455
Stats:
 -hp: 40
 -attack: 40
 -defense: 35
 -special-attack: 50
 -special-defense: 100
 -speed: 70
Types:
 - water
 - poison
 ```

 ### Random encounters

 When you explore there is a chance you encounter a wild pokemon 

 ```
Pokedex > explore canalave-city-area
-  tentacool
-  tentacruel
-  staryu
-  magikarp
-  gyarados
-  wingull
-  pelipper
-  shellos
-  gastrodon
-  finneon
-  lumineon

A wild shellos appears!
Choose an action:
- escape
- battle
- catch
 ```

 Now you are in front of a choice: try to catch it, fight it or just run away. 

 ```
Pokedex/Encounter > battle
Choose a pokemon to fight with shellos
- Enter 'inspect shellos' if you have already catch it to check its stats
- Enter 'pokedex' to check your pokedex
- Enter 'inspect <pokemon-name>' to check the stats of one of your pokemons
- Enter 'choose <pokemon-name>' to start the battle with the chosen pokemon

Pokedex/Battle > pokedex
- gyarados

Pokedex/Battle > choose gyarados
gyarados attacks shellos for 77 damage points
shellos is stunned and got catched!
 ```

 Be aware! If you choose to fight (assuming you have some pokemon to battle with) you even might loose you pokemon of choice.

 ### Different prompts

 When you have a random encounter or a fight you enter in a special kind of prompt in which only the listed commands are possible

 ```
 Pokedex/Encounter >
 ```

 ```
 Pokedex/Battle >
 ```

 Enter the command `exit` to exit this prompts and go back.

 ### Saving and loading progress

 You can save your pokedex to disk with the command `save`, which will produce a directory in which your progress are saved as JSON files.

 In order to load them back just use the `load` command and choose the file you want to load. The files are named as `day_month_year_hour-minute.json` depending on when you saved them.
