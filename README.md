# chessbot2

## modules
### selenium
https://www.zenrows.com/blog/selenium-golang
https://github.com/tebeka/selenium

### robotgo 
https://github.com/go-vgo/robotgo
https://pkg.go.dev/github.com/go-vgo/robotgo

### chess   
https://pkg.go.dev/github.com/notnil/chess#section-readme
https://chess.stackexchange.com/questions/29860/is-there-a-list-of-approximate-elo-ratings-for-each-stockfish-level

## mouse coordinates
```
document.onclick=function(event) {
    var x = event.screenX ;
    var y = event.screenY;
    console.log(x, y) 
}
```

# Lichess
## result
### game in progress
```
result = document.getElementsByClassName("result");
result.length -> 0
```
### game finished
```
result = document.getElementsByClassName("result");
result[0].innerHTML 
```

# Chesscom
## Time selector
```
document.getElementsByClassName("time-selector-button-button")
```
## Board
### Orientation
```
c = document.getElementsByClassName("coordinate-light") 
c[0].textContent    "8" -> play with white
                    "1" -> play with black
```
### Coordinates
```
cl = document.getElementsByClassName("coordinate-light")
cd = document.getElementsByClassName("coordinate-dark")
```
## move list
```
move_list = document.getElementsByClassName("main-line-row")
move = move_list[0].getElementsByClassName("main-line-ply")
```

## clocks
used to check if playing with white or black
```
c = document.getElementsByClassName("clock-time-monospace")
```

## Players
```
p = document.getElementsByClassName("player-avatar")

p[0].innerHTML
'<img alt="Guest2593657151" src="https://www.chess.com/bundles/web/images/white_400.png" width="40" height="40"> <div class="presence-square-component"><div class="presence-square-square" style="width: 1rem; height: 1rem;"></div></div> '

p[1].innerHTML
'<img alt="Guest6542920377" src="https://www.chess.com/bundles/web/images/black_400.png" width="40" height="40"> <div class="presence-square-component"><div class="presence-square-square" style="width: 1rem; height: 1rem;"></div></div> '
```