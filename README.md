# chessbot 2

## selenium
https://www.zenrows.com/blog/selenium-golang
https://github.com/tebeka/selenium

## modules
robotgo https://github.com/go-vgo/robotgo
        https://pkg.go.dev/github.com/go-vgo/robotgo

chess   https://pkg.go.dev/github.com/notnil/chess#section-readme
https://chess.stackexchange.com/questions/29860/is-there-a-list-of-approximate-elo-ratings-for-each-stockfish-level

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



## wait for webelement

```
wait for the element "your turn" for about 60 seconds
this does only work, when there is no time control, which is usualy not the case
err := driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
	yourTurn, _ := driver.FindElement(selenium.ByXPATH, `//*[@id="main-wrap"]/main/div[1]/div[8]/div`)
	if yourTurn != nil {
		return yourTurn.IsDisplayed()
	}
	return false, nil
}, 60*time.Second)
if err != nil {
	return false, err
}
yourTurn, err := driver.FindElement(selenium.ByXPATH, `//*[@id="main-wrap"]/main/div[1]/div[8]/div`)
if err != nil {
	return false, err
}
yt, err := yourTurn.Text()
if err != nil {
	return false, nil
}
fmt.Println("yourturn: ", yt)
return true, nil
```

# chess.com
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