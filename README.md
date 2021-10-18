# gh-screensaver

_being a gh extension that runs animated terminal "screensavers"_

## installation

```
gh extension install vilmibm/gh-screensaver
```

## usage

- `gh screensaver` run a random screensaver
- `gh screensaver -s pipes` run a screensaver by name
- `gh screensaver -l` list available screensavers

Extra configuration options can be passed after a `--`; for example:

```
gh screensaver -smarquee -- --message="hello world" --font="script"
```

## savers

### fireworks

watch a fireworks display.

![fwork2](https://user-images.githubusercontent.com/98482/134737299-aa306b69-ceb4-49c1-95c8-3582d195250c.gif)

`--color` `full` or `off`. Default `full`

### starfield

fly through space.

![starfield](https://user-images.githubusercontent.com/98482/134737341-701d0e7d-476f-4a29-8309-d34b4935c6a3.gif)

`--density` Default `250`. The number of stars to render.
`--speed` Default `4`. Higher is faster.

### pipes

2d pipes draw across the screen.

![pipes2](https://user-images.githubusercontent.com/98482/134737439-34967494-7742-4c55-b92c-da17d6f9f5a9.gif)

`--color` `full` or `off`. Default `full`

### pollock

paint splotches cover the screen.

![pollock](https://user-images.githubusercontent.com/98482/134737473-b5a6a046-58e2-4471-b3c6-3ee191a47af6.gif)

### life
game of life.  

![r-pentomino](https://i.imgur.com/Qq3c0N1.gif)
![dragon](https://media.giphy.com/media/PwIywr183ioixLHqHX/giphy.gif)

`--seed` `glider`,`noise`,`R`,`dragon`,`gun`,or `pulsar`. Default random.  
`--color` `full` or `off`. Default `full`
## author

nate smith <vilmibm@github.com>
