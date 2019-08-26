
Powered by [Ebiten](https://ebiten.org/)

## Todo list:
- ~~Click event should be handled after pressing rather than releasing.~~
- ~~Support web browser.~~
- ~~The grid initially clicked should not response subsequent drag operation.~~
- ~~Generate puzzle randomly with specific size.~~
- ~~Submit & Restart.~~
- ~~Responsive web design.~~
- ~~Add timing.~~
- Add touch support.
- Restart when playing.
- Ensure the uniqueness of answer, the diagonal symmetry of any part of the given table could ruin it.(since nonogram is NP-Complete [paper page29](http://liacs.leidenuniv.nl/assets/2012-01JanvanRijn.pdf), brutal force is an option.)


## Run on web browser
1. Use [Gopherjs](https://github.com/gopherjs/gopherjs) to compile to javascript.
2. Embed `.js` in `.html`