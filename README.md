# Lucky Feet

Entry for [Ebitengine's Holiday Hack](https://itch.io/jam/ebitengine-holiday-hack-2023) (December 23rd 2023 to January 15th 2024). Can be played directly on the browser at [tinne26.github.io/luckyfeet](https://tinne26.github.io/luckyfeet) or downloaded through [github releases](https://github.com/tinne26/luckyfeet/releases/tag/v0.0.2). It's also available through [itch.io](https://tinne26.itch.io/luckyfeet).

![cover](https://github.com/tinne26/luckyfeet/assets/95440833/d6a26142-1348-4f53-ab6a-6a35b561aad4)

I had trouble with motivation, so I only tried to get something working. I wanted to approach the game making a few design choices intentionally contrary to what I'd normally do... and it's been bad. Single screen platformer, lively tone, flat pixel art, what you see is what it is, etc. It's not bad as "what a terrible game", but the feeling was "I don't want to keep developing this any more".

One might wonder why I'm so dumb, but I train really hard for it.

# Controls

Gamepad or keyboard (WASD + SPACE + IOP). The controls are more detailed within the game itself, search for it on the menus.

There's a tic-tac mechanic (see parkour). If you are on the main layer (light brown), you can tic-tac on the back layer (gray). If you are on the front layer (dark brown), you can tic-tac on the main layer. You can't go through walls on the same layer, but can go in front/behind other layers.

# Known Issues

- Little or no optimization. I also decided to double TPS for better input response, which adds insult to injury.
- If you try to load invalid data from clipboard, the game will panic directly instead of showing a nice error. This includes trying the option without knowing what it does. I didn't bother to make the editing nice for players.
- There are many ways to get stuck and trigger behavior that looks outright wrong or buggy.... and some things are kinda rough and unpolished. Yeah.

# License

Code is licensed under the MIT license. Assets are licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).

# Acknowledgements

- Thanks to [kettek](https://github.com/kettek) for organizing the game jam.
- As always, thanks to Hajime Hoshi for creating and developing Ebitengine.

# More

Need more? Check some [extra levels](https://github.com/tinne26/luckyfeet/tree/main/levels).
