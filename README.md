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
## Board orientation
```
c = document.getElementsByClassName("coordinate-light") 
c[0].textContent    "8" -> play with white
                    "1" -> play with black
```
## move list
```
move_list = document.getElementsByClassName("main-line-row")
move = move_list[0].getElementsByClassName("main-line-ply")
```