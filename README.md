# chessbot 2

## selenium
https://www.zenrows.com/blog/selenium-golang
https://github.com/tebeka/selenium

## modules
robotgo https://github.com/go-vgo/robotgo
        https://pkg.go.dev/github.com/go-vgo/robotgo

chess   https://pkg.go.dev/github.com/notnil/chess#section-readme


## mouse coordinates
document.onclick=function(event) {
    var x = event.screenX ;
    var y = event.screenY;
    console.log(x, y) 
}


## result
### game in progress
result = document.getElementsByClassName("result");
result.length -> 0
### game finished
result = document.getElementsByClassName("result");
result[0].innerHTML 
