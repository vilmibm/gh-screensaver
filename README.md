# gh-screensaver

_being a gh extension that runs animated terminal "screensavers"_

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

`--color` `full` or `off`. Default `full`

### starfield

fly through space.

`--density` Default `250`. The number of stars to render.
`--speed` Default `4`. Higher is faster.

### pipes

2d pipes draw across the screen.

`--color` `full` or `off`. Default `full`

### pollock

paint splotches cover the screen.


## author

nate smith <vilmibm@github.com>
